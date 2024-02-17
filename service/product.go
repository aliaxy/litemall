// Package service 相关服务
package service

import (
	"litemall/model"
	"litemall/repository"
)

// IProductService 对于商品服务的接口
type IProductService interface {
	GetProductByID(int64) (*model.Product, error)
	GetAllProduct() ([]*model.Product, error)
	DeleteProductByID(int64) bool
	InsertProduct(*model.Product) (int64, error)
	UpdateProduct(*model.Product) error
}

// ProductService 商品服务实例
type ProductService struct {
	productRepository repository.IProduct
}

// NewProductService 新建服务实例
func NewProductService(repository repository.IProduct) IProductService {
	return &ProductService{
		productRepository: repository,
	}
}

// GetProductByID 根据 ID 查询商品
func (p *ProductService) GetProductByID(id int64) (*model.Product, error) {
	return p.productRepository.SelectByKey(id)
}

// GetAllProduct 查询所有商品
func (p *ProductService) GetAllProduct() ([]*model.Product, error) {
	return p.productRepository.SelectAll()
}

// DeleteProductByID 通过 ID 删除商品
func (p *ProductService) DeleteProductByID(id int64) bool {
	return p.productRepository.Delete(id)
}

// InsertProduct 插入商品
func (p *ProductService) InsertProduct(product *model.Product) (int64, error) {
	return p.productRepository.Insert(product)
}

// UpdateProduct 更新商品
func (p *ProductService) UpdateProduct(product *model.Product) error {
	return p.productRepository.Update(product)
}
