package service

import (
	"errors"
	"fmt"
	"template/common"
	"template/model"

	"gorm.io/gorm"
)

type Task struct{}

func (t Task) Add(task model.Task) (model.Task, error) {
	if err := model.DB.Model(&model.Task{}).Create(&task).Error; err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (t Task) GetInfo(taskID int) (model.Task, error) {
	var task model.Task
	if err := model.DB.Model(&model.Task{}).
		Preload("Poster", func(db *gorm.DB) *gorm.DB {
			return db.Omit("password")
		}).
		Preload("Assignee", func(db *gorm.DB) *gorm.DB {
			return db.Omit("password")
		}).
		Where("id = ?", taskID).
		First(&task).Error; err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (t Task) Delete(taskID int, userID int64) error {
	var record model.Task
	if err := model.DB.Model(&model.Task{}).Where("id = ?", taskID).Select("id, poster_id").First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(fmt.Errorf("未查询到该锅单"), common.OpErr)
		}
		return common.ErrNew(fmt.Errorf("查询锅单失败：%v", err), common.SysErr)
	}
	// 仅发布者可删除锅单，如果没有发布者...说明数据库被入侵了，删了算了
	if record.PosterID != nil && *record.PosterID != userID {
		return common.ErrNew(fmt.Errorf("无权限删除该锅单"), common.AuthErr)
	}
	if err := model.DB.Model(&model.Task{}).Where("id = ?", taskID).Delete(&model.Task{}).Error; err != nil {
		return common.ErrNew(fmt.Errorf("删除锅单失败：%v", err), common.SysErr)
	}

	return nil
}
