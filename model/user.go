package model

type User struct {
	BaseModel
	Name     string `gorm:"column:name;unique;comment:'昵称'" json:"name"`
	Avatar   string `gorm:"column:avatar;comment:'头像'" json:"avatar"`
	Password string `gorm:"column:password;comment:'密码'" json:"-"`
}
