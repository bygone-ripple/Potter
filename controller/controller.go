package controller

import "github.com/gin-gonic/gin"

type Controller struct {
	Auth
	User
	Task
	Comment
}

func New() *Controller {
	Controller := &Controller{}
	return Controller
}

// Upload 上传图片、视频等文件到本地，存储在 uploads 目录下
func (*Controller) Upload(c *gin.Context) {

}
