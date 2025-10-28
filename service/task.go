package service

import (
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
