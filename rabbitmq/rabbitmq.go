// Package rabbitmq 实现了一系列消息队列方法
package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"litemall/model"
	"litemall/service"

	"github.com/streadway/amqp"
)

// MQURL rabbitmq 地址
const MQURL = "amqp://aliaxy:aliaxy@localhost:5672/imooc"

// RabbitMQ 实例
type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string // 队列名称
	Exchange  string // 交换机
	Key       string // key
	MQUrl     string // 连接信息
	sync.Mutex
}

// NewRabbitMQ return *RabbitMQ
func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
	return &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		MQUrl:     MQURL,
	}
}

// Destroy 断开 channel 和 connection
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

// failOnErr 错误处理
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

// NewRabbitMQSimple 创建简单模式实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "")

	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MQUrl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")

	return rabbitmq
}

// PublishSimple 简单模式生产
func (r *RabbitMQ) PublishSimple(message string) error {
	r.Lock()
	defer r.Unlock()
	// 1. 申请队列，如果队列不存在会自动创建，如果存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // 是否持久化
		false, // 是否为自动删除
		false, // 是否具有排他性
		false, // 是否阻塞
		nil,   // 额外属性
	)
	r.failOnErr(err, "failed to declare a queue")

	// 2. 发送消息到队列中
	_ = r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false, // if true, 根据 exchange 类型和 routerkey 规则, 如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false, // if true, 当 exchange 发送消息队列后发现队列上没有绑定消费者, 则会把消息返还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	return nil
}

// ConsumeSimple 简单模式消费
func (r *RabbitMQ) ConsumeSimple(orderService service.IOrderService, productService service.IProductService) {
	// 1. 申请队列，如果队列不存在会自动创建，如果存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName, // 队列名称
		false,       // 是否持久化
		false,       // 是否为自动删除
		false,       // 是否具有排他性
		false,       // 是否阻塞
		nil,         // 额外属性
	)
	if err != nil {
		fmt.Println(err)
	}

	// 消费者流控
	r.channel.Qos(
		1,     // 当前消费者一次能接受的最大消息数量
		0,     // 服务器传递的最大容量（以八位字节为单位）
		false, // 如果设置为true 对channel可用
	)

	// 2. 接收消息
	msg, err := r.channel.Consume(
		r.QueueName, // 队列名称
		"",          // 消费者 用来区分多个消费者
		true,        // 是否自动应答
		false,       // 是否具有排他性
		false,       // if true, 表示不能将同一个 conn 中发送的消息传递给这个 conn 中的消费者
		false,       // 是否阻塞 false 为阻塞
		nil,         // 额外的属性
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	// 启用协程处理消息
	go func() {
		for d := range msg {
			// 消息逻辑处理，可以自行设计逻辑
			log.Printf("Received a message: %s", d.Body)
			message := &model.Message{}
			err := json.Unmarshal([]byte(d.Body), message)
			if err != nil {
				fmt.Println(err)
			}
			// 插入订单
			_, err = orderService.InsertOrderByMessage(message)
			if err != nil {
				fmt.Println(err)
			}

			// 扣除商品数量
			err = productService.SubNumberOne(message.ProductID)
			if err != nil {
				fmt.Println(err)
			}
			// 如果为true表示确认所有未确认的消息，
			// 为false表示确认当前消息
			d.Ack(false)
		}
	}()

	log.Println("[*] Waiting for messages, To exit press CTRL+C")
	<-forever
}
