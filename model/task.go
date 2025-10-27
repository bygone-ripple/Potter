package model

import "time"

// Depart 表示参与的部门，从高到低位分别为技术、视频、美工，如 0b101 表示技术与美工参与
// Level 表示锅单评级，最低为 1，最高为 5
// Status 表示锅单完成状态，1 表示未接取，2 表示进行中，3 表示已完成，4 表示已废止
type Task struct {
	BaseModel
	Name           string      `gorm:"column:name;unique;not null;comment:'锅单名称'" json:"name"`
	Depart         int         `gorm:"column:depart;comment:'参与部门'" json:"depart"`
	Description    string      `gorm:"column:description;comment:'锅单介绍'" json:"decription"`
	Deadline       time.Time   `gorm:"column:ddl;comment:'截止时间'" json:"ddl"`
	Level          int         `gorm:"column:level;default:0;comment:'难度评级'" json:"level"`
	Status         int         `gorm:"column:status;default:0;comment:'锅单状态'" json:"status"`
	CriticalPoints []TimePoint `gorm:"column:critical_points;serializer:json;comment:'关键时间节点'" json:"criticalPoints"`
	Uris           []string    `gorm:"column:uris;serializer:json;comment:'附件资源路径'" json:"uris"`

	Comments   []Comment `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PosterID   *int64    `gorm:"column:poster_id"`
	Poster     User      `gorm:"foreignKey:PosterID;references:ID;constraint:OnDelete:SET NULL;"`
	AssigneeID *int64    `gorm:"column:assignee_id"`
	Assignee   User      `gorm:"foreignKey:AssigneeID;references:ID;constraint:OnDelete:SET NULL"`
}

// CriticalPoints 直接以 json 文本存储在数据库
// URIs 直接以 json 文本存储在数据库
// Task 与 Comment 为 has many
// Task 与 Poster, Assignee 均为 belong to

type TimePoint struct {
	Event string    `json:"event"`
	Time  time.Time `json:"time"`
}

// Depart 数据库值int -> 各部门是否参与的map
func (t *Task) DepartToMap() map[string]bool {
	result := make(map[string]bool)

	result["tech"] = (t.Depart & (1 << 2)) != 0
	result["video"] = (t.Depart & (1 << 1)) != 0
	result["art"] = (t.Depart & (1 << 0)) != 0

	return result
}

func DepartToInt(depart string) int {
	switch depart {
	case "tech":
		return 1 << 2
	case "video":
		return 1 << 1
	case "art":
		return 1 << 0
	}
	return 0
}

func DepartToStr(depart int) string {
	switch depart {
	case 1 << 2:
		return "tech"
	case 1 << 1:
		return "video"
	case 1 << 0:
		return "art"
	}
	return ""
}
