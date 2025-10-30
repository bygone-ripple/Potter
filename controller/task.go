package controller

import (
	"fmt"
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
	task.Status = int(model.NotTaken)

	taskInfo, err := srv.Task.Add(task)
	// 这里返回的 err 没有被包装过
	if err != nil {
		c.Error(common.ErrNew(fmt.Errorf("添加锅单错误：%v", err), common.SysErr))
		return
	}
	c.JSON(http.StatusCreated, ResponseNew(c, taskInfo))
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
	c.JSON(http.StatusOK, ResponseNew(c, struct {
		Total int64      `json:"total"`
		List  []model.Task `json:"list"`
	}{
		Total: total,
		List:  tasks,
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
	// 这里返回的 err 没有被包装过
	if err != nil {
		c.Error(common.ErrNew(fmt.Errorf("获取锅单信息错误：%v", err), common.SysErr))
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

// UpdateInfo 修改锅单信息，仅发布者或接锅人可修改，若要修改发布者和接锅人不应调用此接口
func (t Task) UpdateInfo(c *gin.Context) {
	var uri struct {
		TaskID int `uri:"taskID" binding:"required,min=1"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	var json struct {
		Name           string            `json:"name" binding:"required"`
		Depart         string            `json:"depart" binding:"required,oneof=tech video art"`
		Description    string            `json:"description" binding:"required"`
		Deadline       time.Time         `json:"ddl" binding:"required"`
		Status         int               `json:"status" binding:"required,oneof=1 2 3 4"`
		Level          int               `json:"level" binding:"required,min=1,max=5"`
		CriticalPoints []model.TimePoint `json:"criticalPoints" binding:"required"`
		Uris           []string          `json:"uris" binding:"required"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	session := SessionGet(c, "user-session").(UserSession)
	var task model.Task
	if err := copier.Copy(&task, &json); err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	task.ID = int64(uri.TaskID)
	if json.Depart != "" {
		task.Depart = model.DepartToInt(json.Depart)
	}

	taskInfo, err := srv.Task.UpdateInfo(task, session.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, taskInfo))
}

// UpdateAssignee 修改锅单的接锅人，仅发布者可调用
func (t Task) UpdateAssignee(c *gin.Context) {
	var uri struct {
		TaskID int `uri:"taskID" binding:"required,min=1"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	var json struct {
		AssigneeID int64 `json:"assigneeID" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	session := SessionGet(c, "user-session").(UserSession)
	taskInfo, err := srv.Task.UpdateAssignee(uri.TaskID, json.AssigneeID, session.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, taskInfo))
}

// AddAssignee 将自己添加为锅单的接锅人，并将状态设置成已接取，自由修改不应调用此接口
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

// DeleteAssignee 将自己从接锅人中删除，并将状态设置为未接取，仅接锅人可调用此接口
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
