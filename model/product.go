// Package model 描述不同的数据模型
package model

// Product 商品模型定义
type Product struct {
	ID     int64  `json:"product_id" sql:"product_id" imooc:"product_id"`
	Name   string `json:"product_name" sql:"product_name" imooc:"product_name"`
	Number int64  `json:"product_number" sql:"product_number" imooc:"product_number"`
	Image  string `json:"product_image" sql:"product_image" imooc:"product_image"`
	URL    string `json:"product_url" sql:"product_url" imooc:"product_url"`
}
