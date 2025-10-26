package controller

import (
	"fmt"
	"net/http"
	"template/common"

	"github.com/gin-gonic/gin"
)

type Auth struct{}

func (*Auth) Login(c *gin.Context) {
	var json struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	userInfo, err := srv.Auth.Login(json.Name, json.Password)
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

func (*Auth) Logout(c *gin.Context) {
	if session := SessionGet(c, "user-session"); session == nil {
		c.Error(common.ErrNew(fmt.Errorf("用户未登录"), common.AuthErr))
		return
	}
	SessionDelete(c, "user-session")
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}
