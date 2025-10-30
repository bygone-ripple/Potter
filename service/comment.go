package service

import (
	"fmt"
	"template/common"
	"template/model"
)

type Comment struct{}

func (c Comment) Delete(commentID int, userID int64) error {
	// 允许评论者本人或对应任务的发布者删除
	res := model.DB.Model(&model.Comment{}).
		Joins("LEFT JOIN tasks ON comments.task_id = tasks.id").
		Where("comments.id = ? AND (comments.poster_id = ? OR tasks.poster_id = ?)", commentID, userID, userID).
		Delete(&model.Comment{})
	if res.Error != nil {
		return common.ErrNew(fmt.Errorf("删除评论失败：%v", res.Error), common.SysErr)
	}
	if res.RowsAffected == 0 {
		return common.ErrNew(fmt.Errorf("该评论不存在或无权限删除"), common.OpErr)
	}
	return nil
}
