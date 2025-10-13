package model

import "time"

type Comment struct {
	BaseModel
	Content string    `gorm:"column:content;comment:'评论内容'" json:"content"`
	Time    time.Time `gorm:"column:time;comment:'发布时间'" json:"time"`
	TaskID  int64
	UserID  int64
}