package service

import (
	"fmt"
	"template/common"
	"template/model"
	"template/utils"
)

type Auth struct{}

func (a Auth) Login(name string, password string) (model.User, error) {
	var user model.User
	if err := model.DB.Where("name = ?", name).First(&user).Error; err != nil {
		return model.User{}, common.ErrNew(fmt.Errorf("查询用户信息失败：%v", err), common.SysErr)
	}
	hash := user.Password
	if !utils.CheckPasswordHash(password, hash) {
		return model.User{}, common.ErrNew(fmt.Errorf("密码错误"), common.AuthErr)
	}
	return user, nil
}
