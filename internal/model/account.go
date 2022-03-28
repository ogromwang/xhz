package model

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	gorm.Model
	Username string
	Password string
}

func (m *Account) TableName() string {
	return "account"
}

type AccountDTO struct {
	Id       uint      `json:"id"`
	Username string    `json:"username" binding:"required"`
	Password string    `json:"password" binding:"required"`
	CreateAt time.Time `json:"createAt"`
}
