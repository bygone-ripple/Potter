package controller

import (
	"net/http"
	"template/common"

	"github.com/gin-gonic/gin"
)

type Comment struct{}

// DeleteComment 删除一条自己的或自己发布的任务下的评论
func (*Comment) Delete(c *gin.Context) {
	var uri struct {
		CommentID int `uri:"commentID" binding:"required,min=1"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	session := SessionGet(c, "user-session").(UserSession)
	if err := srv.Comment.Delete(uri.CommentID, session.ID); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
