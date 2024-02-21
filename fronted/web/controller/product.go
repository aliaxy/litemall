package controller

import (
	"html/template"
	"os"
	"path/filepath"
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

var (
	htmlOutPath  = "./fronted/web/htmlProductShow/" // 生成的 html 保存目录
	templatePath = "./fronted/web/view/template/"   // 静态文件模版目录
)

// GetGenerateHtml 生成文件
func (p *ProductController) GetGenerateHtml() {
	// 获取模版文件地址
	contentTmpl, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	productID, err := p.Ctx.URLParamInt64("productID")
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// 获取 html 生成路径
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")

	// 获取模板渲染数据
	product, err := p.ProductService.GetProductByID(productID)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// 生成静态文件
	generateStaticHTML(p.Ctx, contentTmpl, fileName, product)
}

// generateStaticHTML 生成静态文件
func generateStaticHTML(ctx iris.Context, template *template.Template, fileName string, product *model.Product) {
	// 判断静态文件是否存在
	if fileExist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Debug(err)
		}
	}

	// 生成静态文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		ctx.Application().Logger().Debug(err)
	}
	defer file.Close()

	template.Execute(file, product)
}

// fileExist 判断文件是否存在
func fileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
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
