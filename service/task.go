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

func (*Task) AddAssignee(taskID int, userID int64) (model.Task, error) {
	res := model.DB.Model(&model.Task{}).
	Where("id = ? AND assignee_id IS NULL", taskID).
	Update("assignee_id", userID)
	if res.Error != nil {
		return model.Task{}, common.ErrNew(fmt.Errorf("添加接锅人失败：%v", res.Error), common.SysErr)
	}
	if res.RowsAffected == 0 {
		return model.Task{}, common.ErrNew(fmt.Errorf("该锅单已有接锅人"), common.OpErr)
	}
	
	var task model.Task
	if err := model.DB.Model(&model.Task{}).
		Preload("Assignee", func(db *gorm.DB) *gorm.DB {
			return db.Omit("password")
		}).
		Preload("Poster", func(db *gorm.DB) *gorm.DB {
			return db.Omit("password")
		}).
		Where("id = ?", taskID).First(&task).Error; err != nil {
		return model.Task{}, common.ErrNew(fmt.Errorf("接锅人成功更新，查询锅单信息失败：%v", err), common.SysErr)
	}
	return task, nil
}

func (*Task) DeleteAssignee(taskID int, userID int64) error {
	res := model.DB.Model(&model.Task{}).
		Where("id = ? AND assignee_id = ?", taskID, userID).
		Update("assignee_id", nil)
	if res.Error != nil {
		return common.ErrNew(fmt.Errorf("删除接锅人失败：%v", res.Error), common.SysErr)
	}
	if res.RowsAffected == 0 {
		return common.ErrNew(fmt.Errorf("该锅单不存在或接锅人不是您"), common.OpErr)
	}
	return nil
}
