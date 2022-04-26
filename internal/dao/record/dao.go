package record

import (
	"fmt"
	"gorm.io/gorm"
	"strings"

	"github.com/sirupsen/logrus"

	"xiaohuazhu/internal/config"
	"xiaohuazhu/internal/model"
)

type Dao struct {
}

func New() *Dao {
	return &Dao{}
}

func (d *Dao) Add(account *model.RecordMoney) (err error) {
	err = config.AllConn.Db.Save(account).Error
	return
}

func (d *Dao) RecordByFriends(param *model.RecordPageParam, currUser *model.AccountDTO) (resp []*model.RecordPageDTO, err error) {
	record := model.RecordMoney{}
	friend := model.AccountFriend{}
	if err = config.AllConn.Db.Debug().
		Select("friend_ids").
		Where("account_id = ?", currUser.Id).
		First(&friend).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		} else {
			logrus.Errorf("[record|RecordByFriends] DB查询异常, %s", err.Error())
			return
		}
	}
	friendArr := []int64(friend.FriendIds)
	friendArr = append(friendArr, int64(currUser.Id))

	db := config.AllConn.Db.Debug().Table(record.TableName()).
		Select(`
			record_money.created_at,
			record_money.id,
			share, 
			money, 
			describe, 
			image,
			acc.id as account_id, 
			username, 
			profile_picture`).
		Where("(share = true And account_id in (?)) OR account_id = ?", friendArr, currUser.Id).
		Joins("left join account acc on acc.id = record_money.account_id").
		Order(fmt.Sprintf("%s.created_at DESC", record.TableName())).
		Offset(int((param.Page - 1) * param.PageSize)).
		Limit(int(param.PageSize) + 1)

	searchText := strings.Trim(param.SearchText, " ")
	if searchText != "" {
		db = db.Where("acc.username like concat(?, '%')", searchText)
	}

	if err = db.Find(&resp).Error; err != nil {
		logrus.Errorf("[record|RecordByFriends] DB查询异常, %s", err.Error())
	}
	return
}

func (d *Dao) RecordByMe(param *model.RecordPageParam, uid uint) (resp []*model.RecordPageDTO, err error) {
	record := model.RecordMoney{}

	db := config.AllConn.Db.Debug().Table(record.TableName()).
		Select(`
			record_money.created_at,
			record_money.id,
			share, 
			money, 
			describe, 
			image,
			acc.id as account_id, 
			username, 
			profile_picture`).
		Joins("left join account acc on acc.id = record_money.account_id").
		Where("acc.id = ?", uid).
		Order(fmt.Sprintf("%s.created_at DESC", record.TableName())).
		Offset(int((param.Page - 1) * param.PageSize)).
		Limit(int(param.PageSize))

	searchText := strings.Trim(param.SearchText, " ")
	if searchText != "" {
		db = db.Where("acc.username like concat(?, '%')", searchText)
	}

	if err = db.Find(&resp).Error; err != nil {
		logrus.Errorf("[record|RecordByFriends] DB查询异常, %s", err.Error())
	}
	return
}

func (d *Dao) Delete(param *model.RecordDeleteParam, currUser *model.AccountDTO) (record *model.RecordMoney, err error) {
	record = &model.RecordMoney{}

	if err = config.AllConn.Db.Debug().Table(record.TableName()).
		Where("account_id = ? AND id = ?", currUser.Id, param.Id).
		First(record).Error; err != nil {
		logrus.Errorf("[record|Delete] DB查询异常, %s", err.Error())
		return
	}

	if err = config.AllConn.Db.Debug().Unscoped().
		Where("account_id = ? AND id = ?", currUser.Id, param.Id).
		Delete(&model.RecordMoney{}).
		Error; err != nil {
		logrus.Errorf("[record|Delete] 删除异常, %s", err.Error())
		return
	}

	return
}
