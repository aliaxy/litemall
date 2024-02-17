package repository

import (
	"database/sql"
	"errors"

	"litemall/common"
	"litemall/model"
)

// IUserRepository 用户模型对应的接口
type IUserRepository interface {
	Conn() error
	Select(string) (*model.User, error)
	Insert(*model.User) (int64, error)
}

// UserManager 用户接口的具体实现
type UserManager struct {
	table   string
	sqlConn *sql.DB
}

// NewUserManager 创建
func NewUserManager(table string, sqlConn *sql.DB) IUserRepository {
	return &UserManager{
		table:   table,
		sqlConn: sqlConn,
	}
}

// Conn 初始化数据库连接
func (u *UserManager) Conn() error {
	if u.sqlConn == nil {
		mysql, err := common.NewMySQLConn()
		if err != nil {
			return err
		}
		u.sqlConn = mysql
	}
	if u.table == "" {
		u.table = "user"
	}
	return nil
}

// Select 根据 username 查询用户
func (u *UserManager) Select(name string) (user *model.User, err error) {
	if name == "" {
		return &model.User{}, errors.New("条件不能为空！")
	}
	if err = u.Conn(); err != nil {
		return &model.User{}, err
	}

	sql := `select *
			from user
			where user_name = ?`
	rows, err := u.sqlConn.Query(sql, name)
	defer rows.Close()
	if err != nil {
		return &model.User{}, err
	}

	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return &model.User{}, errors.New("用户不存在！")
	}

	user = &model.User{}
	common.DataToStructByTagSQL(result, user)
	return
}

// Insert 插入用户
func (u *UserManager) Insert(user *model.User) (id int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}

	sql := `insert
			into user
			(user_nickname, user_name, user_password)
			values
			(?, ?, ?)`

	stmt, err := u.sqlConn.Prepare(sql)
	if err != nil {
		return id, err
	}

	result, err := stmt.Exec(user.Nickname, user.Name, user.Password)
	if err != nil {
		return id, err
	}

	return result.LastInsertId()
}

// SelectByID 通过 ID 查询
func (u *UserManager) SelectByID(id int64) (user *model.User, err error) {
	if err = u.Conn(); err != nil {
		return &model.User{}, err
	}

	sql := `select *
			from user
			where user_id = ?`

	row, err := u.sqlConn.Query(sql, id)
	if err != nil {
		return &model.User{}, err
	}

	result := common.GetResultRow(row)

	if len(result) == 0 {
		return &model.User{}, errors.New("用户不存在！")
	}

	user = &model.User{}
	common.DataToStructByTagSQL(result, user)
	return
}
