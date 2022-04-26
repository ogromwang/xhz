package model

import (
	"time"

	"gorm.io/gorm"
)

type RecordMoney struct {
	gorm.Model
	AccountId uint
	Share     bool
	Money     float64
	Describe  string
	Image     string
}

func (m *RecordMoney) TableName() string {
	return "record_money"
}

type RecordMoneyDTO struct {
	Share    bool    `form:"share"`
	Money    float64 `form:"money" binding:"required"`
	Describe string  `form:"describe" binding:"required"`
}

type RecordPageParam struct {
	Page       uint   `json:"page" form:"page" binding:"required"`
	PageSize   uint   `json:"pageSize" form:"pageSize" binding:"required"`
	SearchText string `json:"searchText" form:"searchText"`
}

type RecordPageDTO struct {
	CreatedAt      time.Time `json:"createdAt"`
	Id             uint      `json:"id"`
	Share          bool      `json:"share"`
	Money          float64   `json:"money"`
	Describe       string    `json:"describe"`
	Image          string    `json:"image"`
	AccountId      uint      `json:"accountId"`
	Username       string    `json:"username"`
	ProfilePicture string    `json:"profilePicture"`
}
