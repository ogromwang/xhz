package account

import (
	"gorm.io/gorm"
	"strings"
	"xiaohuazhu/internal/config"
	"xiaohuazhu/internal/model"

	"github.com/sirupsen/logrus"
)

type Dao struct {
}

func New() *Dao {
	return &Dao{}
}

func (d *Dao) PageAccount(currId uint, param *model.AccountFriendPageParam) (resp []*model.Account, err error) {
	resp = make([]*model.Account, 0)

	offset := (param.Page - 1) / param.PageSize
	db := config.AllConn.Db.
		Offset(int(offset)).
		Limit(int(param.PageSize)).
		Order("created_at DESC").
		Where("id != ?", currId)

	searchName := strings.Trim(param.Username, " ")
	if searchName != "" {
		db.Where("username like concat(?, '%')", searchName)
	}

	err = db.Find(&resp).Error
	if err != nil {
		logrus.Errorf("[account|PageAccount] 分页查找account, err: [%+v]", err)
	}
	return
}

func (d *Dao) ListFriend(accountId uint) (resp []*model.Account, err error) {
	accountFriend := model.AccountFriend{}
	err = config.AllConn.Db.
		Where("account_id = ?", accountId).
		Select("friend_ids").
		First(&accountFriend).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			resp = make([]*model.Account, 0)
			return
		}
		logrus.Errorf("[account|ListFriend] 发生错误, %s", err.Error())
		return
	}
	friendIds := []int64(accountFriend.FriendIds)

	resp, err = d.List(friendIds)
	return
}

func (d *Dao) List(accountIds []int64) (resp []*model.Account, err error) {
	db := config.AllConn.Db
	if accountIds != nil {
		db.Where("id in (?)", accountIds)
	}
	if err = db.Find(&resp).Error; err != nil {
		logrus.Errorf("[account|List] 发生错误, %s", err.Error())
		return
	}
	return
}

func (d *Dao) Add(account *model.Account) (err error) {
	err = config.AllConn.Db.Save(account).Error
	return
}

func (d *Dao) GetByUsernameOrId(username string, id uint, onlyExist bool) (po *model.Account, err error) {
	po = &model.Account{}
	db := config.AllConn.Db
	if username != "" {
		db.Where("username = ?", username)
	}
	if id != 0 {
		db.Where("id = ?", id)
	}
	if onlyExist {
		db.Select("id")
	}
	err = db.Limit(1).First(po).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		err = nil
		po = nil
		return
	}
	return
}

// ApplyAddFriend 好友申请
func (d *Dao) ApplyAddFriend(acceptFriendId uint, yourId uint) (err error) {
	applyAccountFriend := model.ApplyAccountFriend{
		AccountId: acceptFriendId,
		FriendId:  yourId,
		Status:    0,
	}
	if err = config.AllConn.Db.Save(&applyAccountFriend).Error; err != nil {
		logrus.Errorf("[account|ApplyAddFriend] 发生错误, %s", err.Error())
		return
	}
	return
}
