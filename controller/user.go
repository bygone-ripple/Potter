package controller

import (
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
	user, err := srv.User.Register(json.Name, json.Password)
	if err != nil {
		c.Error(err)
		return
	}
	SessionSet(c, "user-session", UserSession{
		ID:       user.ID,
		UserName: user.Name,
		Level:    1,
	})
	c.JSON(http.StatusOK, ResponseNew(c, user))
}

func (*User) GetInfo(c *gin.Context) {

}

func (*User) UpdateInfo(c *gin.Context) {

}

// GetPostedTasks 获取该用户发布的所有锅单
func (*User) GetPostedTasks(c *gin.Context) {

}

// GetAssignedTasks 获取该用户接取的所有锅单
func (*User) GetAssignedTasks(c *gin.Context) {

}
