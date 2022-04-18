package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Goal struct {
	gorm.Model
	AccountIds pq.Int64Array `json:"account_ids" gorm:"type:int[]"`
	Leader     uint          `json:"leader"`
	Money      float32       `json:"money"`
	Type       int           `json:"type"`
}

func (m *Goal) TableName() string {
	return "goal"
}

type GoalSetDTO struct {
	Id    uint    `json:"id"`
	Money float32 `json:"money" binding:"required"`
	Type  int     `json:"type" binding:"required"`
}

type GoalGetDTO struct {
	Id         uint
	CurrMoney  float32
	TotalMoney float32
	Type       int
}
