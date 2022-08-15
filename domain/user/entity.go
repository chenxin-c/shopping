package user

import (
	"gorm.io/gorm"
)

// 用户模型
type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(30)"`  //名称
	Password  string `gorm:"type:varchar(100)"` //密码1
	Password2 string `gorm:"-"`                 //密码2
	Salt      string `gorm:"type:varchar(100)"` //随机生成
	Token     string `gorm:"type:varchar(500)"`
	IsDeleted bool   //是否删除
	IsAdmin   bool   //是否管理者
}

// 新建用户实例
func NewUser(username, password, password2 string) *User {

	return &User{
		Username:  username,
		Password:  password,
		Password2: password2,
		IsDeleted: false,
		IsAdmin:   false,
	}
}
