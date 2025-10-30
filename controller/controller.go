package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"template/common"
	"time"

	"github.com/gin-gonic/gin"
)

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
	file, err := c.FormFile("file")
	if err != nil {
		c.Error(common.ErrNew(fmt.Errorf("获取文件失败：%v", err), common.ParamErr))
		return
	}
	// 新文件名：原文件名_时间戳.扩展名
	ext := filepath.Ext(file.Filename)
	name := file.Filename[:len(file.Filename)-len(ext)]
	newFilename := fmt.Sprintf("%s_%d%s", name, time.Now().Unix(), ext)

	savePath := filepath.Join("uploads", newFilename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.Error(common.ErrNew(fmt.Errorf("保存文件失败：%v", err), common.SysErr))
		return
	}

	uri := fmt.Sprintf("/static/%s", newFilename)
	c.JSON(http.StatusOK, ResponseNew(c, gin.H{
		"uri": uri,
	}))
}
