package service

import (
	"errors"

	"litemall/model"
	"litemall/repository"

	"golang.org/x/crypto/bcrypt"
)

// IUserService 对于用户服务的接口
type IUserService interface {
	IsPwdSuccess(string, string) (*model.User, bool)
	AddUser(*model.User) (int64, error)
}

// UserService 用户服务实例
type UserService struct {
	UserRepository repository.IUserRepository
}

// NewUserService 新建服务实例
func NewUserService(repository repository.IUserRepository) IUserService {
	return &UserService{
		UserRepository: repository,
	}
}

// IsPwdSuccess 判断密码
func (u *UserService) IsPwdSuccess(username string, password string) (user *model.User, ok bool) {
	user, err := u.UserRepository.Select(username)
	if err != nil {
		return
	}

	ok, _ = validatePassword(password, user.Password)

	if !ok {
		return &model.User{}, false
	}
	return
}

// AddUser 添加用户
func (u *UserService) AddUser(user *model.User) (id int64, err error) {
	pwdByte, err := generatePassword(user.Password)
	if err != nil {
		return
	}
	user.Password = string(pwdByte)
	return u.UserRepository.Insert(user)
}

// generatePassword 生成密码
func generatePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// validatePassword 验证密码
func validatePassword(password string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false, errors.New("密码比对错误！")
	}
	return true, nil
}
