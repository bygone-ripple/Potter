package model

import "time"

type Comment struct {
	BaseModel
	Content  string    `gorm:"column:content;comment:'评论内容'" json:"content"`
	Time     time.Time `gorm:"column:time;comment:'发布时间'" json:"time"`

	TaskID   int64     `gorm:"column:task_id" json:"taskID"`
	PosterID int64     `gorm:"column:poster_id" json:"posterID"`
	Poster   *User      `gorm:"foreignKey:PosterID;references:ID;constraint:OnDelete:CASCADE" json:"poster,omitempty"`
}
// Comment 与 Poster 为 belong to