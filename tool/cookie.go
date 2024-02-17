// Package tool 工具包
package tool

import (
	"net/http"

	"github.com/kataras/iris/v12"
)

// GlobalCookie 设置全局cookie
func GlobalCookie(ctx iris.Context, name, value string) {
	ctx.SetCookie(&http.Cookie{Name: name, Value: value, Path: "/"})
}
