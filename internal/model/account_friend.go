package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AccountFriend struct {
	gorm.Model
	AccountId uint          `json:"account_id"`
	FriendIds pq.Int64Array `json:"friend_ids" gorm:"type:int[]"`
}

func (m *AccountFriend) TableName() string {
	return "account_friend"
}

type ApplyAccountFriend struct {
	gorm.Model
	AccountId uint `json:"account_id"`
	FriendId  uint `json:"friend_id"`
	Status    int
}

func (m *ApplyAccountFriend) TableName() string {
	return "account_friend_apply"
}

type AccountFriendPageParam struct {
	Page     uint   `json:"page" binding:"required"`
	PageSize uint   `json:"pageSize" binding:"required"`
	Username string `json:"username"`
}

type ApplyAddFriendParam struct {
	Id uint `json:"id" binding:"required"`
}
