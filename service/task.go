package service

import (
	"fmt"
	"template/common"
	"template/model"
	"time"

	"gorm.io/gorm"
)

type Task struct{}

func (t Task) Add(task model.Task) (model.Task, error) {
	if err := model.DB.Model(&model.Task{}).Create(&task).Error; err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (t Task) Get(pager common.PagerForm, name string, depart string, status int, level int) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	db := model.DB.Model(&model.Task{})
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	filter := model.Task{
		Depart: model.DepartToInt(depart),
		Status: status,
		Level:  level,
	}
	db = db.Where(&filter)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, common.ErrNew(fmt.Errorf("获取锅单总数失败：%v", err), common.SysErr)
	}
	if err := db.Offset((pager.Page - 1) * pager.Limit).Limit(pager.Limit).Find(&tasks).Error; err != nil {
		return nil, 0, common.ErrNew(fmt.Errorf("获取锅单失败：%v", err), common.SysErr)
	}
	return tasks, total, nil
}

func (t Task) GetInfo(taskID int) (model.Task, error) {
	var task model.Task
	if err := model.DB.Model(&model.Task{}).
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Poster").Order("time DESC")
		}).
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
	var task model.Task
	if err := model.DB.Model(&model.Task{}).Where("id = ?", taskID).First(&task).Error; err != nil {
		return common.ErrNew(fmt.Errorf("查询锅单失败：%v", err), common.SysErr)
	}
	if task.PosterID != nil && *task.PosterID != userID {
		return common.ErrNew(fmt.Errorf("您不是该锅单的负责人，无法删除"), common.OpErr)
	}
	if err := model.DB.Model(&model.Task{}).Delete(&task).Error; err != nil {
		return common.ErrNew(fmt.Errorf("删除锅单失败：%v", err), common.SysErr)
	}

	return nil
}

func (t Task) UpdateInfo(taskInfo model.Task, userID int64) (model.Task, error) {
	res := model.DB.Model(&model.Task{}).Omit("poster_id", "assignee_id").Where("id = ? AND (poster_id = ? OR assignee_id = ?)", taskInfo.ID, userID, userID).Updates(&taskInfo)
	if res.Error != nil {
		return model.Task{}, common.ErrNew(fmt.Errorf("更新锅单失败：%v", res.Error), common.SysErr)
	}
	if res.RowsAffected == 0 {
		return model.Task{}, common.ErrNew(fmt.Errorf("该锅单不存在或无权限更新"), common.OpErr)
	}
	task, err := t.GetInfo(int(taskInfo.ID))
	if err != nil {
		return model.Task{}, common.ErrNew(fmt.Errorf("更新锅单成功，查询锅单信息失败：%v", err), common.SysErr)
	}
	return task, nil
}

func (t Task) UpdateAssignee(taskID int, assigneeID int64, userID int64) (model.Task, error) {
	res := model.DB.Model(&model.Task{}).Where("id = ? AND poster_id = ?", taskID, userID).Updates(map[string]any{
		"assignee_id": assigneeID,
	})
	if res.Error != nil {
		return model.Task{}, common.ErrNew(fmt.Errorf("更新接锅人失败：%v", res.Error), common.SysErr)
	}
	if res.RowsAffected == 0 {
		return model.Task{}, common.ErrNew(fmt.Errorf("该锅单不存在或发布者不是您"), common.OpErr)
	}
	task, err := t.GetInfo(taskID)
	if err != nil {
		return model.Task{}, common.ErrNew(fmt.Errorf("更新接锅人成功，查询锅单信息失败：%v", err), common.SysErr)
	}
	return task, nil
}

func (t Task) AddAssignee(taskID int, userID int64) (model.Task, error) {
	res := model.DB.Model(&model.Task{}).
		Where("id = ? AND assignee_id IS NULL", taskID).
		Updates(map[string]any{
			"assignee_id": userID,
			"status":      int(model.InProgress),
		})
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

func (t Task) DeleteAssignee(taskID int, userID int64) error {
	res := model.DB.Model(&model.Task{}).
		Where("id = ? AND assignee_id = ?", taskID, userID).
		Updates(map[string]any{
			"assignee_id": nil,
			"status":      int(model.NotTaken),
		})
	if res.Error != nil {
		return common.ErrNew(fmt.Errorf("删除接锅人失败：%v", res.Error), common.SysErr)
	}
	if res.RowsAffected == 0 {
		return common.ErrNew(fmt.Errorf("该锅单不存在或接锅人不是您"), common.OpErr)
	}
	return nil
}

func (t Task) PostComment(taskID int, posterID int64, content string) (model.Comment, error) {
	comment := model.Comment{
		Content:  content,
		Time:     time.Now(),
		TaskID:   int64(taskID),
		PosterID: posterID,
	}
	if err := model.DB.Model(&model.Comment{}).Create(&comment).Error; err != nil {
		return model.Comment{}, common.ErrNew(fmt.Errorf("发布评论失败：%v", err), common.SysErr)
	}
	comment.Poster = &model.User{}
	if err := model.DB.Model(&model.User{}).Where("id = ?", posterID).First(comment.Poster).Error; err != nil {
		return model.Comment{}, common.ErrNew(fmt.Errorf("发布评论成功，查询评论者信息失败：%v", err), common.SysErr)
	}
	return comment, nil
}
