package controller

import (
	"net/http"
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
	c.JSON(http.StatusCreated, ResponseNew(c, struct {
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
func (t Task) Get(c *gin.Context) {
	var query struct {
		common.PagerForm
		Name   string `form:"name" binding:"omitempty"`
		Depart string `form:"depart" binding:"omitempty,oneof=tech video art"`
		Status int    `form:"status" binding:"omitempty,oneof=1 2 3 4"`
		Level  int    `form:"level" binding:"omitempty,oneof=1 2 3 4 5"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	tasks, total, err := srv.Task.Get(query.PagerForm, query.Name, query.Depart, query.Status, query.Level)
	if err != nil {
		c.Error(err)
		return
	}
	type taskInfo struct {
		ID          int64     `json:"id"`
		Name        string    `json:"name"`
		Depart      string    `json:"depart"`
		Description string    `json:"description"`
		Deadline    time.Time `json:"ddl"`
		Level       int       `json:"level"`
	}
	var responseData []taskInfo
	for _, task := range tasks {
		responseData = append(responseData, taskInfo{
			ID:          task.ID,
			Name:        task.Name,
			Depart:      model.DepartToStr(task.Depart),
			Description: task.Description,
			Deadline:    task.Deadline,
			Level:       task.Level,
		})
	}
	c.JSON(http.StatusOK, ResponseNew(c, struct {
		Total int64      `json:"total"`
		List  []taskInfo `json:"list"`
	}{
		Total: total,
		List:  responseData,
	}))
}

// GetInfo 通过 ID 获取锅单详细信息
func (t Task) GetInfo(c *gin.Context) {
	var uri struct {
		TaskID int `uri:"taskID" binding:"required,min=1"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	taskInfo, err := srv.Task.GetInfo(uri.TaskID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, taskInfo))
}

// Delete 删除锅单，仅发布者可删除
func (t Task) Delete(c *gin.Context) {
	var uri struct {
		TaskID int `uri:"taskID" binding:"required,min=1"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	session := SessionGet(c, "user-session").(UserSession)
	if err := srv.Task.Delete(uri.TaskID, session.ID); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusNoContent, ResponseNew(c, nil))
}

func (t Task) UpdateInfo(c *gin.Context) {

}

// AddAssignee 将自己添加为锅单的接锅人，自由修改不应调用此接口
func (t Task) AddAssignee(c *gin.Context) {
	var uri struct {
		TaskID int `uri:"taskID" binding:"required,min=1"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	session := SessionGet(c, "user-session").(UserSession)
	taskInfo, err := srv.Task.AddAssignee(uri.TaskID, session.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, ResponseNew(c, taskInfo))
}

// DeleteAssignee 将自己从接锅人中删除
func (t Task) DeleteAssignee(c *gin.Context) {
	var uri struct {
		TaskID int `uri:"taskID" binding:"required,min=1"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	session := SessionGet(c, "user-session").(UserSession)
	err := srv.Task.DeleteAssignee(uri.TaskID, session.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusNoContent, ResponseNew(c, nil))
}

func (t Task) PostComment(c *gin.Context) {
	var uri struct {
		TaskID int `uri:"taskID" binding:"required,min=1"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	var json struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}

	session := SessionGet(c, "user-session").(UserSession)
	commentInfo, err := srv.Task.PostComment(uri.TaskID, session.ID, json.Content)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, ResponseNew(c, commentInfo))
}
