// Package main 程序入口
package main

import (
	"github.com/kataras/iris/v12"
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

	// 注册控制器

	// 启动服务
	app.Run(
		iris.Addr(":8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
