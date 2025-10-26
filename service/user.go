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

	if err := model.DB.Model(&model.User{}).Create(&user).Error; err != nil {
		return model.User{}, common.ErrNew(fmt.Errorf("创建用户记录失败：%v", err), common.SysErr)
	}

	var userInfo model.User
	if err := model.DB.Model(&model.User{}).
		Select("id, name").Where("name = ?", name).First(&userInfo).Error; err != nil {
		return model.User{}, common.ErrNew(fmt.Errorf("获取用户信息失败：%v", err), common.SysErr)
	}
	return userInfo, nil
}
