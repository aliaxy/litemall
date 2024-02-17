package model

// User 用户模型定义
type User struct {
	ID       int64  `json:"user_id" sql:"user_id" imooc:"user_id"`
	Nickname string `json:"user_nickname" sql:"user_nickname" imooc:"user_nickname"`
	Name     string `json:"user_name" sql:"user_name" imooc:"user_name"`
	Password string `json:"user_password" sql:"user_password" imooc:"user_password"`
}
