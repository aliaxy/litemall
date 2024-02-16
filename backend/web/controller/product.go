// Package controller 商品控制层
package controller

import (
	"strconv"

	"litemall/common"
	"litemall/model"
	"litemall/service"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// ProductController 商品对外控制
type ProductController struct {
	Ctx            iris.Context
	ProductService *service.ProductService
}

// GetList 获取商品列表
func (p *ProductController) GetList() mvc.View {
	productList, _ := p.ProductService.GetAllProduct()
	return mvc.View{
		Name: "product/view.html",
		Data: iris.Map{
			"productList": productList,
		},
	}
}

// PostUpdate 修改商品
func (p *ProductController) PostUpdate() {
	product := new(model.Product)
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{
		TagName: "imooc",
	})
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	err := p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	p.Ctx.Redirect("/product/list")
}

// GetAdd 添加商品
func (p *ProductController) GetAdd() mvc.View {
	return mvc.View{
		Name: "product/add.html",
	}
}

// PostAdd 添加商品
func (p *ProductController) PostAdd() {
	product := new(model.Product)
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{
		TagName: "imooc",
	})
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	_, err := p.ProductService.InsertProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	p.Ctx.Redirect("/product/list")
}

// GetManager 管理商品
func (p *ProductController) GetManager() mvc.View {
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(id)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name: "product/manager.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

// GetDelete 删除商品
func (p *ProductController) GetDelete() {
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	ok := p.ProductService.DeleteProductByID(id)
	if ok {
		p.Ctx.Application().Logger().Debug("删除商品成功, id: ", id)
	} else {
		p.Ctx.Application().Logger().Debug("删除商品失败, id: ", id)
	}
}
