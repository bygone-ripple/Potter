package service

import (
	"fmt"
	"template/common"
	"template/model"
	"template/utils"
)

type User struct{}

func (u User) Register(name string, password string) (model.User, error) {
	encrypted, encryptErr := utils.Encrypt(password)
	if encryptErr != nil {
		return model.User{}, common.ErrNew(fmt.Errorf("加密密码时出错：%v", encryptErr), common.SysErr)
	}
	user := model.User{
		Name:     name,
		Password: encrypted,
	}

	var count int64
	if err := model.DB.Model(&User{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return model.User{}, common.ErrNew(fmt.Errorf("查询昵称是否重复失败：%v", err), common.SysErr)
	}

	if count > 0 {
		return model.User{}, common.ErrNew(fmt.Errorf("该昵称已被其他人使用"), common.OpErr)
	}
	// TODO: Error 1364 (HY000): Field 'created_at' doesn't have a default value
	if err := model.DB.Model(&User{}).Create(&user).Error; err != nil {
		return model.User{}, common.ErrNew(fmt.Errorf("创建用户记录失败：%v", err), common.SysErr)
	}
	return user, nil
}
