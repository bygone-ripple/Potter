package controller

import "github.com/gin-gonic/gin"

type Task struct{}

func (*Task) Add(c *gin.Context) {

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
