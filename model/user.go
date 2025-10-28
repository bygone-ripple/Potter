package model

type User struct {
	BaseModel
	Name     string `gorm:"column:name;comment:'昵称'" json:"name"`
	Password string `gorm:"column:password;comment:'密码'" json:"-"`
}
