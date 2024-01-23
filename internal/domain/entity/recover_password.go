package entity

import (
	"time"

	"gorm.io/gorm"
)

type RecoverPassword struct {
	Model
	Token    string `gorm:"column:token;type:varchar(500);not null"`
	Nickname string `gorm:"column:nickname;type:varchar(255);not null"`
	State    bool   `gorm:"column:state;type:tinyint(1);not null"`
}

func (m RecoverPassword) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return
}

func (m RecoverPassword) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}
