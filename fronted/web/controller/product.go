package controller

import (
	"strconv"

	"litemall/model"
	"litemall/service"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

// ProductController 商品详情控制
type ProductController struct {
	Ctx            iris.Context
	ProductService service.IProductService
	OrderService   service.IOrderService
	Session        *sessions.Session
}

// GetDetail 商品详情页面
func (p *ProductController) GetDetail() mvc.View {
	product, err := p.ProductService.GetProductByID(1)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

// GetOrder 订单页面
func (p *ProductController) GetOrder() mvc.View {
	productID, err := p.Ctx.URLParamInt("productID")
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	userString := p.Ctx.GetCookie("uid")
	userID, err := strconv.ParseInt(userString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	var orderID int64
	showMessage := "抢购失败"
	// 判断商品数量是否满足需求
	if product.Number > 0 {
		// 扣除商品数量
		product.Number--
		p.ProductService.UpdateProduct(product)

		// 创建订单
		order := &model.Order{
			UserID:    userID,
			ProductID: int64(productID),
			Status:    model.OrderSuccess,
		}
		orderID, err = p.OrderService.InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		} else {
			showMessage = "抢购成功"
		}
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderID":     orderID,
			"showMessage": showMessage,
		},
	}
}
