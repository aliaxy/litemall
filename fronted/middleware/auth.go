// Package middleware 各种中间件
package middleware

import "github.com/kataras/iris/v12"

// AuthConProduct 检查是否登录
func AuthConProduct(ctx iris.Context) {
	// 从 cookie 中得到用户 id
	uid := ctx.GetCookie("uid")

	if uid == "" {
		ctx.Application().Logger().Debug("必须先登录")
		ctx.Redirect("/user/login")
		return
	}

	ctx.Application().Logger().Debug("已经登录")
	ctx.Next()
}
