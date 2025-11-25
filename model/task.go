package model

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Depart 表示参与的部门，从高到低位分别为技术、视频、美工，如 0b101 表示技术与美工参与
// Level 表示锅单评级，最低为 1，最高为 5
// Status 表示锅单完成状态，1 表示未接取，2 表示进行中，3 表示已完成，4 表示已废止
type Task struct {
	BaseModel
	Name        string    `gorm:"column:name;unique;not null;comment:'锅单名称'" json:"name"`
	Depart      int       `gorm:"column:depart;comment:'参与部门'" json:"-"`
	Description string    `gorm:"column:description;comment:'锅单介绍'" json:"description"`
	Deadline    time.Time `gorm:"column:ddl;comment:'截止时间'" json:"ddl"`
	Level       int       `gorm:"column:level;default:0;comment:'难度评级'" json:"level"`
	Status      int       `gorm:"column:status;default:0;comment:'锅单状态'" json:"status"`
	// CriticalPoints []TimePoint `gorm:"column:critical_points;serializer:json;comment:'关键时间节点'" json:"criticalPoints"`
	Uris []string `gorm:"column:uris;serializer:json;comment:'附件资源路径'" json:"uris"`

	Comments   []Comment `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"comments,omitempty"`
	PosterID   *int64    `gorm:"column:poster_id" json:"posterID,omitempty"`
	Poster     *User     `gorm:"foreignKey:PosterID;references:ID;constraint:OnDelete:SET NULL;" json:"poster,omitempty"`
	AssigneeID *int64    `gorm:"column:assignee_id" json:"assigneeID,omitempty"`
	Assignee   *User     `gorm:"foreignKey:AssigneeID;references:ID;constraint:OnDelete:SET NULL" json:"assignee,omitempty"`
}

// CriticalPoints 直接以 json 文本存储在数据库
// URIs 直接以 json 文本存储在数据库
// Task 与 Comment 为 has many
// Task 与 Poster, Assignee 均为 belong to

// type TimePoint struct {
// 	Event string    `json:"event"`
// 	Time  time.Time `json:"time"`
// }

type Status int

const (
	NotTaken   Status = 1 // 未接取
	InProgress Status = 2 // 进行中
	Completed  Status = 3 // 已完成
	Abandoned  Status = 4 // 已废止
)

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

func (t Task) BeforeDelete(tx *gorm.DB) (err error) {
	if t.ID == 0 {
		return nil
	}
	newName := fmt.Sprintf("%s__deleted_%d", t.Name, time.Now().Unix())
	return tx.Model(&Task{}).Where("id = ?", t.ID).Update("name", newName).Error
}

// MarshalJSON 自定义 Task 结构体的 JSON 编码，以便将 Depart 字段转换为字符串
func (t Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID          int64     `json:"id,omitempty"`
		Name        string    `json:"name"`
		Depart      string    `json:"depart"`
		Description string    `json:"description"`
		Deadline    time.Time `json:"ddl"`
		Level       int       `json:"level"`
		Status      int       `json:"status"`
		// CriticalPoints []TimePoint `json:"criticalPoints"`
		Uris       []string  `json:"uris"`
		Comments   []Comment `json:"comments,omitempty"`
		PosterID   *int64    `json:"posterID,omitempty"`
		Poster     *User     `json:"poster,omitempty"`
		AssigneeID *int64    `json:"assigneeID,omitempty"`
		Assignee   *User     `json:"assignee,omitempty"`
	}{
		ID:          t.ID,
		Name:        t.Name,
		Depart:      DepartToStr(t.Depart),
		Description: t.Description,
		Deadline:    t.Deadline,
		Level:       t.Level,
		Status:      t.Status,
		// CriticalPoints: t.CriticalPoints,
		Uris:       t.Uris,
		Comments:   t.Comments,
		PosterID:   t.PosterID,
		Poster:     t.Poster,
		AssigneeID: t.AssigneeID,
		Assignee:   t.Assignee,
	})
}

// UnmarshalJSON 实现自定义 JSON 解码，以便将 Depart 字段从字符串转换为整数
func (t *Task) UnmarshalJSON(data []byte) error {
	type Alias Task
	aux := &struct {
		*Alias
		Depart string `json:"depart"`
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Depart != "" {
		t.Depart = DepartToInt(aux.Depart)
	}
	return nil
}
