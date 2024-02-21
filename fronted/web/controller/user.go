// Package controller 前台控制层
package controller

import (
	"fmt"
	"strconv"

	"litemall/encrypt"
	"litemall/model"
	"litemall/service"
	"litemall/tool"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

// UserController 用户控制层
type UserController struct {
	Ctx     iris.Context
	Service service.IUserService
	Session *sessions.Session
}

// GetRegister 注册页面
func (u *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",
	}
}

// PostRegister 发送注册请求
func (u *UserController) PostRegister() {
	var (
		nickname = u.Ctx.FormValue("user_nickname")
		username = u.Ctx.FormValue("user_name")
		password = u.Ctx.FormValue("user_password")
	)
	// ozzo-validation
	user := &model.User{
		Name:     username,
		Nickname: nickname,
		Password: password,
	}

	_, err := u.Service.AddUser(user)
	u.Ctx.Application().Logger().Debug(err)
	if err != nil {
		u.Ctx.Redirect("/user/error")
		return
	}
	u.Ctx.Redirect("/user/login")
	return
}

// GetLogin 登录页面
func (u *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

// PostLogin 登录请求
func (u *UserController) PostLogin() mvc.Response {
	// 1.获取用户提交的表单信息
	username := u.Ctx.FormValue("user_name")
	password := u.Ctx.FormValue("user_password")

	// 2、验证账号密码正确
	user, ok := u.Service.IsPwdSuccess(username, password)
	if !ok {
		return mvc.Response{
			Path: "/user/login",
		}
	}

	// 3、写入用户ID到cookie中
	tool.GlobalCookie(u.Ctx, "uid", strconv.FormatInt(user.ID, 10))
	uidByte := []byte(strconv.FormatInt(user.ID, 10))
	uidString, err := encrypt.EnPasswordCode(uidByte)
	if err != nil {
		fmt.Println(err)
	}

	// 4、写入用户 ID 到浏览器
	tool.GlobalCookie(u.Ctx, "sign", uidString)

	return mvc.Response{
		Path: "/product/",
	}
}
