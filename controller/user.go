package controller

import "github.com/gin-gonic/gin"

type User struct{}

func (*User) Register(c *gin.Context) {

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
