// Package main 程序入口
package main

import (
	"context"
	"log"

	"litemall/backend/web/controller"
	"litemall/common"
	"litemall/repository"
	"litemall/service"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	// 创建 iris 实例
	app := iris.New()

	// 设置错误模式
	app.Logger().SetLevel("debug")

	// 注册模版
	template := iris.HTML("./backend/web/view", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)

	// 设置模版目标
	app.HandleDir("/assets", "./backend/web/assets")

	// 出现异常跳转指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	// 连接数据库
	db, err := common.NewMySQLConn()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册控制器
	productRepository := repository.NewProductManager("product", db)
	productSerivce := service.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productSerivce)
	product.Handle(new(controller.ProductController))

	orderRepository := repository.NewOrderManager("order", db)
	orderService := service.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controller.OrderController))

	// 启动服务
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
