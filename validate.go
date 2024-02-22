// Package main 分布式架构
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"litemall/common"
	"litemall/encrypt"
	"litemall/model"
	"litemall/rabbitmq"
)

var (
	hosts     = []string{"http://127.0.0.1", "http://127.0.0.1"}
	localhost = "http://127.0.0.1"
	port      = "8081"
	// GetOneIP 数量控制接口服务器内网IP，或者getone的SLB内网IP
	GetOneIP = "127.0.0.1"
	// GetOnePort 对应端口
	GetOnePort       = "8084"
	hashConsistent   *common.Consistent
	rabbitMQValidate *rabbitmq.RabbitMQ
	accessControl    = &AccessControl{
		sourcesArray: make(map[int]interface{}),
	}
)

// AccessControl 访问控制
type AccessControl struct {
	// 用来存放用户想要存放的信息
	sourcesArray map[int]interface{}
	sync.RWMutex
}

// GetNewRecord 获取制定的数据
func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.sourcesArray[uid]
	return data
}

// SetNewRecord 设置记录
func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	m.sourcesArray[uid] = "hello litemall"
	m.RWMutex.Unlock()
}

// GetDistributedRight 得到顺时针分配的
func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	// 获取用户UID
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}

	// 采用一致性hash算法，根据用户ID，判断获取具体机器
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	// 判断是否为本机
	if hostRequest == localhost {
		// 执行本机数据读取和校验
		return m.GetDataFromMap(uid.Value)
	}
	// 不是本机充当代理访问数据返回结果
	return GetDataFromOtherMap(hostRequest, req)
}

// GetDataFromMap 获取本机map
// 并且处理业务逻辑
// 返回的结果类型为bool类型
func (m *AccessControl) GetDataFromMap(uid string) bool {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := m.GetNewRecord(uidInt)
	// 执行逻辑判断
	// 测试使用
	if data != nil {
		return true
	}
	return true
}

// GetDataFromOtherMap 获取其它节点处理结果
func GetDataFromOtherMap(host string, request *http.Request) bool {
	hostURL := "http://" + host + ":" + port + "/checkRight"
	response, body, err := GetCurl(hostURL, request)
	if err != nil {
		return false
	}
	// 判断状态
	if response.StatusCode == 200 {
		return string(body) == "true"
	}
	return false
}

// GetCurl 模拟请求
func GetCurl(hostURL string, request *http.Request) (response *http.Response, body []byte, err error) {
	// 获取Uid
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}
	// 获取sign
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}

	// 模拟接口访问，
	client := &http.Client{}
	req, err := http.NewRequest("GET", hostURL, nil)
	if err != nil {
		return
	}

	// 手动指定，排查多余cookies
	cookieUID := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	// 添加cookie到模拟的请求中
	req.AddCookie(cookieUID)
	req.AddCookie(cookieSign)

	// 获取返回结果
	response, err = client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return
	}
	body, err = io.ReadAll(response.Body)
	return
}

// CheckRight 检测
func CheckRight(w http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

// Check 执行正常业务逻辑
func Check(w http.ResponseWriter, r *http.Request) {
	// 执行正常业务逻辑
	fmt.Println("执行check！")
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"]) <= 0 {
		w.Write([]byte("false"))
		return
	}
	productString := queryForm["productID"][0]
	fmt.Println(productString)
	// 获取用户cookie
	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false"))
		return
	}

	// 1.分布式权限验证
	right := accessControl.GetDistributedRight(r)
	if right == false {
		w.Write([]byte("false"))
		return
	}
	// 2.获取数量控制权限，防止秒杀出现超卖现象
	hostURL := "http://" + GetOneIP + ":" + GetOnePort + "/getOne"
	responseValidate, validateBody, err := GetCurl(hostURL, r)
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	// 判断数量控制接口请求状态
	if responseValidate.StatusCode == 200 {
		if string(validateBody) == "true" {
			// 整合下单
			// 1.获取商品ID
			productID, err := strconv.ParseInt(productString, 10, 64)
			if err != nil {

				w.Write([]byte("false"))
				return
			}
			// 2.获取用户ID
			userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {

				w.Write([]byte("false"))
				return
			}

			// 3.创建消息体
			message := model.NewMessage(userID, productID)
			// 类型转化
			byteMessage, err := json.Marshal(message)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			// 4.生产消息
			err = rabbitMQValidate.PublishSimple(string(byteMessage))
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("false"))
	return
}

// Auth 统一验证拦截器
// 每个接口都需要验证
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证！")
	// 添加基于cookie的权限验证
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

// CheckUserInfo 验证用户信息
func CheckUserInfo(r *http.Request) error {
	// 获取 cookie
	uid, err := r.Cookie("uid")
	if err != nil {
		return errors.New("用户 uid 获取失败")
	}

	// 获取用户加密串
	sign, err := r.Cookie("sign")
	if err != nil {
		return errors.New("用户加密串获取失败")
	}

	// 对信息进行解密
	signByte, err := encrypt.DePasswordCode(sign.Value)
	if err != nil {
		return errors.New("用户加密串已被篡改")
	}

	if checkInfo(uid.Value, string(signByte)) {
		return nil
	}
	return errors.New("身份校验失败")
}

// checkInfo 自定义逻辑判断
func checkInfo(checkStr, signStr string) bool {
	return checkStr == signStr
}

func main() {
	// 负载均衡器设置
	// 使用一致性哈希算法
	hashConsistent = common.NewConsistent()
	// 添加节点
	for _, v := range hosts {
		hashConsistent.Add(v)
	}

	localIP, err := common.GetIntranceIP()
	if err != nil {
		fmt.Println(err)
	}

	localhost = localIP
	fmt.Println(localhost)

	rabbitMQValidate = rabbitmq.NewRabbitMQSimple("imoocProduct")
	defer rabbitMQValidate.Destroy()

	// 过滤器
	filter := common.NewFilter()
	//@TODO 优化注册拦截器
	filter.RegisterFilterURI("/check", Auth)
	filter.RegisterFilterURI("/checkRight", Auth)
	// 2、启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight", filter.Handle(CheckRight))
	// 启动服务
	http.ListenAndServe(":8083", nil)
}
