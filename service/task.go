package service

import (
	"template/model"
)

type Task struct{}

func (t Task) Add(task model.Task) (model.Task, error) {
	if err := model.DB.Model(&model.Task{}).Create(&task).Error; err != nil {
		return model.Task{}, err
	}
	return task, nil
}
