// Package main 分布式架构
package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"

	"litemall/common"
	"litemall/encrypt"
)

var (
	hosts          = []string{"http://127.0.0.1", "http://127.0.0.1"}
	localhost      = "http://127.0.0.1"
	port           = "8081"
	hashConsistent *common.Consistent
	accessControl  = &AccessControl{
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
	if data != nil {
		return true
	}
	return false
}

// GetDataFromOtherMap 获取其它节点处理结果
func GetDataFromOtherMap(host string, request *http.Request) bool {
	// 获取Uid
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return false
	}
	// 获取sign
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return false
	}

	// 模拟接口访问，
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+host+":"+port+"/check", nil)
	if err != nil {
		return false
	}

	// 手动指定，排查多余cookies
	cookieUID := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	// 添加cookie到模拟的请求中
	req.AddCookie(cookieUID)
	req.AddCookie(cookieSign)

	// 获取返回结果
	response, err := client.Do(req)
	if err != nil {
		return false
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false
	}

	// 判断状态
	if response.StatusCode == http.StatusOK {
		return string(body) == "true"
	}
	return false
}

// Auth 统一验证拦截器
// 每个接口都需要验证
func Auth(rw http.ResponseWriter, req *http.Request) (err error) {
	fmt.Println("执行 auth")
	// 添加基于 cookie 的权限验证
	err = CheckUserInfo(req)
	return
}

// Check 执行正常业务逻辑
func Check(rw http.ResponseWriter, req *http.Request) error {
	fmt.Println("执行 check")
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
	// 过滤器
	filter := common.NewFilter()
	// 注册拦截器
	filter.RegisterFilterURI("/check", Auth)
	// 启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	http.ListenAndServe(":8082", nil)
}
