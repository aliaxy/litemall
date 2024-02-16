// Package common 常用工具
package common

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 导入但不使用 init
)

// NewMySQLConn 创建 MySQL 连接
func NewMySQLConn() (db *sql.DB, err error) {
	dsn := "root:211010@tcp(127.0.0.1:13306)/litemall?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = sql.Open("mysql", dsn)
	return
}

// GetResultRow 获取返回值,取一条
func GetResultRow(rows *sql.Rows) map[string]string {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]string)
	for rows.Next() {
		// 将行数据保存到record字典
		rows.Scan(scanArgs...)
		for i, v := range values {
			if v != nil {
				var value string
				switch v.(type) {
				case []byte:
					value = string(v.([]byte))
				case int64:
					value = fmt.Sprintf("%d", v)
				default:
					value = fmt.Sprintf("%v", v)
				}
				record[columns[i]] = value
			}
		}

	}
	return record
}

// GetResultRows 获取所有
func GetResultRows(rows *sql.Rows) map[int]map[string]string {
	columns, _ := rows.Columns()
	vals := make([][]byte, len(columns))
	scans := make([]interface{}, len(columns))
	for k := range vals {
		scans[k] = &vals[k]
	}

	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		rows.Scan(scans...)
		row := make(map[string]string)
		for k, v := range vals {
			key := columns[k]
			row[key] = string(v)
		}
		result[i] = row
		i++
	}
	return result
}
