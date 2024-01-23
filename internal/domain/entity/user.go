package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Model
	Name     string `gorm:"column:name;type:varchar(255);not null"`
	Surname  string `gorm:"column:surname;varchar(255);not null"`
	Email    string `gorm:"column:email;type:varchar(255);not null"`
	NickName string `gorm:"column:nick_name;type:varchar(45);not null"`
	IsAdmin  bool   `gorm:"column:is_admin;type:tinyint(1);default:0"`
	State    bool   `gorm:"column:state;type:tinyint(1);not null"`
	Password string `gorm:"column:password;varchar(255);not null"`
}

func (m User) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return
}

func (m User) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}
