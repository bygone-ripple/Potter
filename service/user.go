package service

import (
	"fmt"
	"template/common"
	"template/model"
	"template/utils"
)

type User struct{}

func (u User) Register(name string, avatar string, password string) (model.User, error) {
	hash, hashErr := utils.HashPassword(password)
	if hashErr != nil {
		return model.User{}, common.ErrNew(fmt.Errorf("哈希密码时出错：%v", hashErr), common.SysErr)
	}
	user := model.User{
		Name:     name,
		Avatar:   avatar,
		Password: hash,
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
	return user, nil
}

func (u User) GetInfo(userID int64) (model.User, error) {
	var user model.User
	if err := model.DB.Model(&model.User{}).Where("id = ?", userID).First(&user).Error; err != nil {
		return model.User{}, common.ErrNew(fmt.Errorf("获取用户信息失败：%v", err), common.SysErr)
	}
	return user, nil
}

func (u User) UpdateInfo(id int64, name string, avatar string, password string) (model.User, error) {
	var user model.User
	if err := model.DB.Model(&model.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return model.User{}, common.ErrNew(fmt.Errorf("获取用户信息失败：%v", err), common.SysErr)
	}

	if name != "" {
		user.Name = name
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if password != "" {
		hash, hashErr := utils.HashPassword(password)
		if hashErr != nil {
			return model.User{}, common.ErrNew(fmt.Errorf("哈希密码时出错：%v", hashErr), common.SysErr)
		}
		user.Password = hash
	}

	if err := model.DB.Save(&user).Error; err != nil {
		return model.User{}, common.ErrNew(fmt.Errorf("更新用户信息失败：%v", err), common.SysErr)
	}

	return user, nil
}

func (u User) GetPostedTasks(userID int64) ([]model.Task, error) {
	var tasks []model.Task
	if err := model.DB.Model(&model.Task{}).Where("poster_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, common.ErrNew(fmt.Errorf("获取用户发布的锅单失败：%v", err), common.SysErr)
	}
	return tasks, nil
}

func (u User) GetAssignedTasks(userID int64) ([]model.Task, error) {
	var tasks []model.Task
	if err := model.DB.Model(&model.Task{}).Where("assignee_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, common.ErrNew(fmt.Errorf("获取用户接取的锅单失败：%v", err), common.SysErr)
	}
	return tasks, nil
}
