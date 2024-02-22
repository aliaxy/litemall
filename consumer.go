package main

import (
	"fmt"

	"litemall/common"
	"litemall/rabbitmq"
	"litemall/repository"
	"litemall/service"
)

func main() {
	db, err := common.NewMySQLConn()
	if err != nil {
		fmt.Println(err)
	}
	// 创建product数据库操作实例
	product := repository.NewProductManager("product", db)
	// 创建product serivce
	productService := service.NewProductService(product)
	// 创建Order数据库实例
	order := repository.NewOrderManager("order", db)
	// 创建order Service
	orderService := service.NewOrderService(order)

	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("imoocProduct")
	rabbitmqConsumeSimple.ConsumeSimple(orderService, productService)
}
