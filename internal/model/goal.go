package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Goal struct {
	gorm.Model
	Name       string        `json:"name"`
	AccountIds pq.Int64Array `json:"account_ids" gorm:"type:int[]"`
	Leader     uint          `json:"leader"`
	Money      float64       `json:"money"`
	Type       int           `json:"type"`
}

func (m *Goal) TableName() string {
	return "goal"
}

type GoalSetDTO struct {
	Id         uint    `json:"id"  binding:"required"`
	Money      float64 `json:"money" binding:"required"`
	AccountIds []int64 `json:"account_ids" binding:"required"`
	Type       int     `json:"type" binding:"required"`
	Name       string  `json:"name"`
}

type GoalCreateDTO struct {
	Name  string  `json:"Name" binding:"required"`
	Money float64 `json:"money" binding:"required"`
	Type  int     `json:"type" binding:"required"`
}

type GoalGetDTO struct {
	Id         uint
	Name       string
	Goal       float64
	CurrMoney  float64
	Type       int
	AccountIds []int64
}
