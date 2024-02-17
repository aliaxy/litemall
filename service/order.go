package service

import (
	"litemall/model"
	"litemall/repository"
)

// IOrderService 对于订单服务的接口
type IOrderService interface {
	GetOrderByID(int64) (*model.Order, error)
	GetAllOrder() ([]*model.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
	DeleteOrderByID(int64) bool
	InsertOrder(*model.Order) (int64, error)
	UpdateOrder(*model.Order) error
}

// OrderService 订单服务实例
type OrderService struct {
	OrderRepository repository.IOrder
}

// NewOrderService 新建服务实例
func NewOrderService(repository repository.IOrder) IOrderService {
	return &OrderService{
		OrderRepository: repository,
	}
}

// GetOrderByID 根据 ID 查询订单
func (o *OrderService) GetOrderByID(id int64) (*model.Order, error) {
	return o.OrderRepository.SelectByKey(id)
}

// GetAllOrder 查询所有订单
func (o *OrderService) GetAllOrder() ([]*model.Order, error) {
	return o.OrderRepository.SelectAll()
}

// GetAllOrderInfo 查询所有订单信息
func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	return o.OrderRepository.SelectAllWithInfo()
}

// DeleteOrderByID 通过 ID 删除订单
func (o *OrderService) DeleteOrderByID(id int64) bool {
	return o.OrderRepository.Delete(id)
}

// InsertOrder 插入订单
func (o *OrderService) InsertOrder(order *model.Order) (int64, error) {
	return o.OrderRepository.Insert(order)
}

// UpdateOrder 更新订单
func (o *OrderService) UpdateOrder(order *model.Order) error {
	return o.OrderRepository.Update(order)
}
