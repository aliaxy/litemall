// Package repository 数据模型对应的接口
package repository

import (
	"database/sql"

	"litemall/common"
	"litemall/model"
)

// IProduct 商品模型对应的接口
type IProduct interface {
	Conn() error
	Insert(*model.Product) (int64, error)
	Delete(int64) bool
	Update(*model.Product) error
	SelectByKey(int64) (*model.Product, error)
	SelectAll() ([]*model.Product, error)
}

// ProductManager 商品接口的具体实现
type ProductManager struct {
	table   string
	sqlConn *sql.DB
}

// NewProductManager 创建
func NewProductManager(table string, sqlConn *sql.DB) IProduct {
	return &ProductManager{
		table:   table,
		sqlConn: sqlConn,
	}
}

// Conn 初始化数据库连接
func (p *ProductManager) Conn() (err error) {
	if p.sqlConn == nil {
		mysql, err := common.NewMySQLConn()
		if err != nil {
			return err
		}
		p.sqlConn = mysql
	}

	if p.table == "" {
		p.table = "product"
	}

	return
}

// Insert 插入
func (p *ProductManager) Insert(product *model.Product) (ID int64, err error) {
	// 判断连接是否存在
	if err = p.Conn(); err != nil {
		return
	}

	// 准备 sql
	sql := `insert
			into product
			(product_name, product_number, product_image, product_url)
			values
			(?, ?, ?, ?)`
	stmt, err := p.sqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}

	// 执行 sql
	result, err := stmt.Exec(product.Name, product.Number, product.Image, product.URL)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// Delete 删除
func (p *ProductManager) Delete(ID int64) bool {
	// 判断连接是否存在
	if err := p.Conn(); err != nil {
		return false
	}

	// 准备 sql
	sql := `delete
			from product
			where product_id = ?`
	stmt, err := p.sqlConn.Prepare(sql)
	if err != nil {
		return false
	}

	// 执行 sql
	_, err = stmt.Exec(ID)
	if err != nil {
		return false
	}

	return true
}

// Update 更新
func (p *ProductManager) Update(product *model.Product) (err error) {
	// 判断连接是否存在
	if err = p.Conn(); err != nil {
		return
	}

	// 准备 sql
	sql := `update product
			set product_name = ?,
				product_number = ?,
				product_image = ?,
				product_url = ?
			where product_id = ?`
	stmt, err := p.sqlConn.Prepare(sql)
	if err != nil {
		return
	}

	// 执行 sql
	_, err = stmt.Exec(product.Name, product.Number, product.Image, product.URL, product.ID)
	if err != nil {
		return
	}

	return
}

// SelectByKey 查询指定 ID 的记录
func (p *ProductManager) SelectByKey(id int64) (product *model.Product, err error) {
	// 判断连接是否存在
	if err = p.Conn(); err != nil {
		return nil, err
	}

	// 准备 sql
	sql := `select *
			from product
			where product_id = ?`

	// 执行 sql
	row, err := p.sqlConn.Query(sql, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return nil, err
	}

	product = new(model.Product)
	common.DataToStructByTagSQL(result, product)
	return
}

// SelectAll 查询所有记录
func (p *ProductManager) SelectAll() (products []*model.Product, err error) {
	// 判断连接是否存在
	if err = p.Conn(); err != nil {
		return nil, err
	}

	// 准备 sql
	sql := `select *
		from product`

	// 执行 sql
	rows, err := p.sqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}

	for _, v := range result {
		product := &model.Product{}
		common.DataToStructByTagSQL(v, product)
		products = append(products, product)
	}

	return
}
