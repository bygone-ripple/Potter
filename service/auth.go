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
	decrypted, decryptErr := utils.Decrypt(user.Password)
	if decryptErr != nil {
		return model.User{}, common.ErrNew(fmt.Errorf("解密密码时出错：%v", decryptErr), common.SysErr)
	}
	if decrypted != password {
		return model.User{}, common.ErrNew(fmt.Errorf("密码错误"), common.AuthErr)
	}
	return user, nil
}
