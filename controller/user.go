package controller

import (
	"fmt"
	"net/http"
	"template/common"

	"github.com/gin-gonic/gin"
)

type User struct{}

func (*User) Register(c *gin.Context) {
	var json struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	userInfo, err := srv.User.Register(json.Name, json.Password)
	if err != nil {
		c.Error(err)
		return
	}
	SessionSet(c, "user-session", UserSession{
		ID:       userInfo.ID,
		UserName: userInfo.Name,
		Level:    1,
	})
	c.JSON(http.StatusOK, ResponseNew(c, struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}{
		ID:   userInfo.ID,
		Name: userInfo.Name,
	}))
}

func (*User) GetInfo(c *gin.Context) {

}

// UpdateInfo 修改用户信息，需要验证用户身份，若传入参数为空则不修改该项
func (*User) UpdateInfo(c *gin.Context) {
	var json struct {
		ID       int64  `json:"id" binding:"required"`
		Name     string `json:"name" binding:"omitempty"`
		Password string `json:"password" binding:"omitempty"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	if session := SessionGet(c, "user-session"); session == nil || session.(UserSession).ID != json.ID {
		c.Error(common.ErrNew(fmt.Errorf("无权限修改该用户信息"), common.AuthErr))
		return
	}
	userInfo, err := srv.User.UpdateInfo(json.ID, json.Name, json.Password)
	if err != nil {
		c.Error(err)
		return
	}
	SessionSet(c, "user-session", UserSession{
		ID:       userInfo.ID,
		UserName: userInfo.Name,
		Level:    1,
	})
	c.JSON(http.StatusOK, ResponseNew(c, struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}{
		ID:   userInfo.ID,
		Name: userInfo.Name,
	}))
}

// GetPostedTasks 获取该用户发布的所有锅单
func (*User) GetPostedTasks(c *gin.Context) {

}

// GetAssignedTasks 获取该用户接取的所有锅单
func (*User) GetAssignedTasks(c *gin.Context) {

}
