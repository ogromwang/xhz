package model

import "gorm.io/gorm"

type RecordMoney struct {
	gorm.Model
	AccountId uint
	Share     bool
	Money     float32
	Describe  string
	Image     string
}

func (m *RecordMoney) TableName() string {
	return "record_money"
}

type RecordMoneyDTO struct {
	AccountId uint    `json:"accountId" binding:"required"`
	Share     bool    `json:"share" binding:"required"`
	Money     float32  `json:"money" binding:"required"`
	Describe  string  `json:"describe"`
	Image     string  `json:"image"`
}
