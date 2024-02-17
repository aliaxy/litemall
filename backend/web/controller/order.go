package controller

import (
	"litemall/service"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// OrderController 订单对外控制
type OrderController struct {
	Ctx          iris.Context
	OrderService *service.OrderService
}

// Get 查询所有订单
func (o *OrderController) Get() mvc.View {
	orderArray, err := o.OrderService.GetAllOrderInfo()
	if err != nil {
		o.Ctx.Application().Logger().Debug("查询订单信息失败")
	}

	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderArray,
		},
	}
}
