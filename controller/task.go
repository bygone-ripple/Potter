package controller

import (
	"template/common"
	"template/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type Task struct{}

// Add 新建锅单，目前部门仅支持选一个，后续可扩展为多选
func (*Task) Add(c *gin.Context) {
	var json struct {
		Name           string            `json:"name" binding:"required"`
		Depart         string            `json:"depart" binding:"required"`
		Description    string            `json:"description" binding:"required"`
		Deadline       time.Time         `json:"ddl" binding:"required"`
		Level          int               `json:"level" binding:"required,min=1,max=5"`
		CriticalPoints []model.TimePoint `json:"criticalPoints" binding:"required"`
		Uris           []string          `json:"uris" binding:"omitempty"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	var task model.Task
	if err := copier.Copy(&task, &json); err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	// 接口中间件有验证登录，故此处 session 不为空
	session := SessionGet(c, "user-session").(UserSession)
	task.PosterID = &session.ID
	task.Depart = model.DepartToInt(json.Depart)

	taskInfo, err := srv.Task.Add(task)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, ResponseNew(c, struct {
		ID             int64             `json:"id"`
		Name           string            `json:"name"`
		Depart         string            `json:"depart"`
		Description    string            `json:"description"`
		Deadline       time.Time         `json:"ddl"`
		Level          int               `json:"level"`
		CriticalPoints []model.TimePoint `json:"criticalPoints"`
		Uris           []string          `json:"uris"`
	}{
		ID:             taskInfo.ID,
		Name:           taskInfo.Name,
		Depart:         json.Depart,
		Description:    taskInfo.Description,
		Deadline:       taskInfo.Deadline,
		Level:          taskInfo.Level,
		CriticalPoints: taskInfo.CriticalPoints,
		Uris:           taskInfo.Uris,
	}))
}

// Get 通过名称、部门等查询参数获取锅单列表，支持分页查询
func (*Task) Get(c *gin.Context) {

}

func (*Task) GetInfo(c *gin.Context) {

}

func (*Task) Delete(c *gin.Context) {

}

func (*Task) UpdateInfo(c *gin.Context) {

}

// AddAssignee 将自己添加为锅单的接锅人
func (*Task) AddAssignee(c *gin.Context) {

}

// DeleteAssignee 将自己从接锅人中删除
func (*Task) DeleteAssignee(c *gin.Context) {

}

func (*Task) PostComment(c *gin.Context) {

}
