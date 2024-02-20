// Package main
package main

import (
	"context"
	"time"

	"litemall/common"
	"litemall/fronted/middleware"
	"litemall/fronted/web/controller"
	"litemall/repository"
	"litemall/service"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

func main() {
	// 创建 iris 实例
	app := iris.New()

	// 设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")

	// 注册模板
	template := iris.HTML("./fronted/web/view", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)

	// 设置模版目标
	app.HandleDir("/public", "./fronted/web/public")
	app.HandleDir("/html", "./fronted/web/htmlProductShow")

	// 出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	// 连接数据库
	db, _ := common.NewMySQLConn()

	sess := sessions.New(sessions.Config{
		Cookie:  "AdminCookie",
		Expires: 600 * time.Minute,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册控制器
	user := repository.NewUserManager("user", db)
	userService := service.NewUserService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx, sess.Start)
	userPro.Handle(new(controller.UserController))

	product := repository.NewProductManager("product", db)
	productService := service.NewProductService(product)
	productPro := mvc.New(app.Party("/product"))
	productPro.Router.Use(middleware.AuthConProduct)
	productPro.Register(productService, ctx, sess.Start)
	productPro.Handle(new(controller.ProductController))

	app.Run(
		iris.Addr("localhost:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
