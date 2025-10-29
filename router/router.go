package router

import (
	"template/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Use(middleware.Error)
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	r.Static("/static", "./uploads")
	apiRouter := r.Group("/api")
	{
		apiRouter.POST("/uploads", ctr.Upload, middleware.CheckRole(1))
		authRouter := apiRouter.Group("/auth")
		{
			authRouter.POST("/login", ctr.Auth.Login)
			authRouter.POST("/logout", ctr.Auth.Logout)
		}
		userRouter := apiRouter.Group("/users")
		{
			userRouter.POST("/", ctr.User.Register)
			selfRouter := userRouter.Group("/me")
			{
				selfRouter.Use(middleware.CheckRole(1))
				selfRouter.PUT("/", ctr.User.UpdateInfo)
				selfRouter.GET("/posted-tasks", ctr.User.GetPostedTasks)
				selfRouter.GET("/assigned-tasks", ctr.User.GetAssignedTasks)
			}
		}
		taskRouter := apiRouter.Group("/tasks")
		{
			taskRouter.Use(middleware.CheckRole(1))
			taskRouter.POST("/", ctr.Task.Add)
			taskRouter.GET("/", ctr.Task.Get)
			taskRouter.GET("/:taskID", ctr.Task.GetInfo)
			taskRouter.DELETE("/:taskID", ctr.Task.Delete)
			taskRouter.PUT("/:taskID", ctr.Task.UpdateInfo)
			taskRouter.POST("/:taskID/assignees/me", ctr.Task.AddAssignee)
			taskRouter.DELETE("/:taskID/assignees/me", ctr.Task.DeleteAssignee)
			taskRouter.POST("/:taskID/comments", ctr.Task.PostComment)
		}
		commentRouter := apiRouter.Group("/comments")
		{
			commentRouter.Use(middleware.CheckRole(1))
			commentRouter.DELETE("/:commentID", ctr.Comment.DeleteComment)
		}
	}
}
