package types

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(20);" json:"username"`
	Password string `gorm:"column:password;type:varchar(100);" json:"password"`
	Number   int64  `gorm:"column:number;type:int;" json:"number"`
	Content  string `gorm:"column:content;type:varchar(200);" json:"content"`
	Company  string `gorm:"column:company;type:varchar(100);" json:"company"`
	UserJob  string `gorm:"column:user_job;type:varchar(100);" json:"user_job"`
	Birthday string `gorm:"column:birthday;type:varchar(20);" json:"birthday"`
	City     string `gorm:"column:city;type:varchar(100);" json:"city"`
	Email    string `gorm:"column:email;type:varchar(100);" json:"email"`
}

func (u User) TableName() string {
	return "user"
}
