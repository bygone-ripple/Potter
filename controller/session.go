package controller

import (
	"encoding/gob"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ID 设置为 User 的 ID，Level 目前均设置成 1
type UserSession struct {
	ID       int64
	UserName string
	Level    int
}

func _SessionSave(ss sessions.Session) {
	if err := ss.Save(); err != nil {
		log.Fatalf("session save error: %v", err)
	}
}

func SessionGet(c *gin.Context, name string) any {
	session := sessions.Default(c)
	return session.Get(name)
}

func SessionSet(c *gin.Context, name string, body any) {
	session := sessions.Default(c)
	if body == nil {
		return
	}
	gob.Register(body)
	session.Set(name, body)
	_SessionSave(session)
}

func SessionUpdate(c *gin.Context, name string, body any) {
	SessionSet(c, name, body)
}

func SessionClear(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	_SessionSave(session)
}

func SessionDelete(c *gin.Context, name string) {
	session := sessions.Default(c)
	session.Delete(name)
	_SessionSave(session)
}
