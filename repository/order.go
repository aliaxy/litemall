package repository

import (
	"database/sql"

	"litemall/common"
	"litemall/model"
)

// IOrder 订单模型对应的接口
type IOrder interface {
	Conn() error
	Insert(*model.Order) (int64, error)
	Delete(int64) bool
	Update(*model.Order) error
	SelectByKey(int64) (*model.Order, error)
	SelectAll() ([]*model.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

// OrderManager 订单接口的具体实现
type OrderManager struct {
	table   string
	sqlConn *sql.DB
}

// NewOrderManager 创建
func NewOrderManager(table string, sqlConn *sql.DB) IOrder {
	return &OrderManager{
		table:   table,
		sqlConn: sqlConn,
	}
}

// Conn 初始化数据库连接
func (o *OrderManager) Conn() error {
	if o.sqlConn == nil {
		mysql, err := common.NewMySQLConn()
		if err != nil {
			return err
		}
		o.sqlConn = mysql
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

// Insert 插入
func (o *OrderManager) Insert(order *model.Order) (id int64, err error) {
	// 判断连接是否存在
	if err = o.Conn(); err != nil {
		return
	}

	// 准备 sql
	sql := `insert
			into order
			(user_id, product_id, order_status)
			values
			(?, ?, ?)`
	stmt, err := o.sqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}

	// 执行 sql
	result, err := stmt.Exec(order.UserID, order.ProductID, order.Status)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// Delete 删除
func (o *OrderManager) Delete(id int64) bool {
	// 判断连接是否存在
	if err := o.Conn(); err != nil {
		return false
	}

	// 准备 sql
	sql := `delete
			from order
			where order_id = ?`
	stmt, err := o.sqlConn.Prepare(sql)
	if err != nil {
		return false
	}

	// 执行 sql
	_, err = stmt.Exec(id)
	if err != nil {
		return false
	}

	return true
}

// Update 更新
func (o *OrderManager) Update(order *model.Order) (err error) {
	// 判断连接是否存在
	if err = o.Conn(); err != nil {
		return
	}

	// 准备 sql
	sql := `update product
			set user_id = ?,
				product_id = ?,
				order_status = ?,
			where order_id = ?`
	stmt, err := o.sqlConn.Prepare(sql)
	if err != nil {
		return
	}

	// 执行 sql
	_, err = stmt.Exec(order.UserID, order.ProductID, order.Status, order.ID)
	if err != nil {
		return
	}

	return
}

// SelectByKey 查询指定 ID 的记录
func (o *OrderManager) SelectByKey(id int64) (order *model.Order, err error) {
	// 判断连接是否存在
	if err = o.Conn(); err != nil {
		return &model.Order{}, err
	}

	// 准备 sql
	sql := `select *
			from order
			where order_id = ?`

	// 执行 sql
	row, err := o.sqlConn.Query(sql, id)
	if err != nil {
		return &model.Order{}, err
	}
	defer row.Close()
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &model.Order{}, err
	}

	order = new(model.Order)
	common.DataToStructByTagSQL(result, order)
	return
}

// SelectAll 查询所有记录
func (o *OrderManager) SelectAll() (orders []*model.Order, err error) {
	// 判断连接是否存在
	if err = o.Conn(); err != nil {
		return nil, err
	}

	// 准备 sql
	sql := `select *
			from order`

	// 执行 sql
	rows, err := o.sqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}

	for _, v := range result {
		order := &model.Order{}
		common.DataToStructByTagSQL(v, order)
		orders = append(orders, order)
	}

	return
}

// SelectAllWithInfo 查询订单所有商品信息
func (o *OrderManager) SelectAllWithInfo() (orderMap map[int]map[string]string, err error) {
	// 判断连接是否存在
	if err = o.Conn(); err != nil {
		return nil, err
	}

	// 准备 sql
	sql := "select o.order_id, p.product_name, o.order_status " +
		"from `order` as o " +
		"join product as p on o.product_id = p.product_id"

	// 执行 sql
	rows, err := o.sqlConn.Query(sql)
	if err != nil {
		return nil, err
	}

	return common.GetResultRows(rows), nil
}
