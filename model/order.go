package model

// Order 商品模型定义
type Order struct {
	ID        int64  `json:"order_id" sql:"order_id" imooc:"order_id"`
	UserID    int64  `json:"user_id" sql:"user_id" imooc:"user_id"`
	ProductID int64  `json:"product_id" sql:"product_id" imooc:"product_id"`
	Status    string `json:"order_status" sql:"order_status" imooc:"order_status"`
}

const (
	OrderWait    = iota // OrderWait 等待
	OrderSuccess        // OrderSuccess 成功
	OrderFailed         // OrderFailed 失败
)
