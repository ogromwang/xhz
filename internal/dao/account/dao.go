package account

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"strings"
	"xiaohuazhu/internal/config"
	"xiaohuazhu/internal/model"
	"xiaohuazhu/internal/util"

	"github.com/sirupsen/logrus"
)

type Dao struct {
}

func New() *Dao {
	return &Dao{}
}

func (d *Dao) UpdatePicture(uid uint, path string) (err error) {
	account := model.Account{}
	err = config.AllConn.Db.
		Model(account).
		Where("id = ?", uid).
		UpdateColumns(map[string]interface{}{"profile_picture": path}).
		Error
	return
}

func (d *Dao) PageFindAccount(currUid uint, param *model.AccountFriendPageParam) (resp []*model.Account, count int64, err error) {
	subQuery := config.AllConn.Db.Table("account_friend").
		Select("regexp_split_to_table(array_to_string(friend_ids, ','), ',')::int").
		Where("account_id = ?", currUid)
	db := config.AllConn.Db.Table("account").
		Where("account.id != ?", currUid).
		Where("account.id not in (?)", subQuery).
		Where("account.username like concat(?, '%')", param.Username)

	if err = db.Count(&count).Error; err != nil {
		logrus.Errorf("[account|PageFindAccount] 分页查找account, err: [%+v]", err)
		return
	}
	if err = db.Offset(int((param.Page - 1) * param.PageSize)).Limit(int(param.PageSize)).Find(&resp).Error; err != nil {
		logrus.Errorf("[account|PageFindAccount] 分页查找account, err: [%+v]", err)
		return
	}
	return
}

func (d *Dao) PageAccount(notIn []uint, param *model.AccountFriendPageParam) (resp []*model.Account, err error) {
	resp = make([]*model.Account, 0)

	offset := (param.Page - 1) * param.PageSize
	db := config.AllConn.Db.
		Offset(int(offset)).
		Limit(int(param.PageSize)).
		Order("created_at DESC")

	searchName := strings.Trim(param.Username, " ")
	if searchName != "" {
		db = db.Where("username like concat(?, '%')", searchName)
	}
	if notIn != nil && len(notIn) > 0 {
		db = db.Where("id not in (?)", notIn)
	}

	err = db.Find(&resp).Error
	if err != nil {
		logrus.Errorf("[account|PageAccount] 分页查找account, err: [%+v]", err)
	}
	return
}

func (d *Dao) PageFriend(currUid uint, param *model.AccountFriendPageParam) (resp []*model.Account, count int64, err error) {
	subQuery := config.AllConn.Db.Table("account_friend").
		Select("regexp_split_to_table(array_to_string(friend_ids, ','), ',')::int").
		Where("account_id = ?", currUid)
	db := config.AllConn.Db.Table("account").
		Where("account.id in (?)", subQuery)

	search := strings.Trim(param.Username, " ")
	if search != "" {
		db = db.Where("account.username like concat(?, '%')", search)
	}
	if err = db.Count(&count).Error; err != nil {
		logrus.Errorf("[account|PageFriend] 发生错误, %s", err.Error())
		return
	}
	if err = db.Offset(int((param.Page - 1) * param.PageSize)).Limit(int(param.PageSize)).Find(&resp).Error; err != nil {
		logrus.Errorf("[account|PageFriend] 分页查找account, err: [%+v]", err)
		return
	}
	return
}

func (d *Dao) List(accountIds []int64) (resp []*model.Account, err error) {
	db := config.AllConn.Db
	if accountIds == nil || len(accountIds) == 0 {
		return
	}
	if err = db.Where("id in (?)", accountIds).Find(&resp).Error; err != nil {
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
		db = db.Where("username = ?", username)
	}
	if id != 0 {
		db = db.Where("id = ?", id)
	}
	if onlyExist {
		db = db.Select("id")
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
	var num int64
	config.AllConn.Db.Model(applyAccountFriend).
		Where("account_id = ? and friend_id = ?", acceptFriendId, yourId).
		Where("status in (?)", []int{0, 1}).Select("id").
		Count(&num)
	if num > 0 {
		logrus.Warnf("[account|ApplyAddFriend] 警告, 已经申请了 acceptFriendId: [%d] yourId: [%d]", acceptFriendId, yourId)
		return
	}

	if err = config.AllConn.Db.Save(&applyAccountFriend).Error; err != nil {
		logrus.Errorf("[account|ApplyAddFriend] 发生错误, %s", err.Error())
		return
	}
	return
}

// HandleAddFriend 处理好友申请
func (d *Dao) HandleAddFriend(friendId uint, yourId uint, status int) (err error) {
	// 接受
	tx := config.AllConn.Db.Begin()

	if err = tx.Debug().
		Model(model.ApplyAccountFriend{}).
		Where("status = 0").
		Where("(account_id = ? and friend_id = ?) or (account_id = ? and friend_id = ?)", yourId, friendId, friendId, yourId).
		UpdateColumns(map[string]interface{}{"status": status}).
		Error; err != nil {
		tx.Rollback()
		logrus.Errorf("[account|HandleAddFriend] 更新 ApplyAccountFriend 发生错误, %s", err.Error())
		return
	}

	// 写入，或修改一条数据
	if status == 1 {
		err = d.handleRelationship(tx, yourId, friendId)
		err = d.handleRelationship(tx, friendId, yourId)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	return
}

func (d *Dao) handleRelationship(tx *gorm.DB, yourId uint, friendId uint) (err error) {
	accountFriend := model.AccountFriend{}
	// 新增
	if err = tx.Debug().Where("account_id = ?", yourId).First(&accountFriend).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			accountFriend.AccountId = yourId
			accountFriend.FriendIds = pq.Int64Array{int64(friendId)}
			if err = tx.Debug().Save(&accountFriend).Error; err != nil {
				logrus.Errorf("[account|HandleAddFriend] 新增发生错误, %s", err.Error())
				return
			}
			err = nil
			return
		}
		// 异常
		logrus.Errorf("[account|HandleAddFriend] 发生错误, %s", err.Error())
		return
	}

	// 修改
	int64s := []int64(accountFriend.FriendIds)
	if index := util.IntContains(int64s, int64(friendId)); index != -1 {
		// 已经添加了
		logrus.Warnf("[account|HandleAddFriend] 已经添加了")
		return
	}

	int64s = append(int64s, int64(friendId))

	newArray := pq.Int64Array(int64s)
	accountFriend.FriendIds = newArray

	if err = tx.Debug().
		Model(accountFriend).
		Where("account_id = ?", yourId).
		UpdateColumns(map[string]interface{}{"friend_ids": accountFriend.FriendIds}).
		Error; err != nil {
		logrus.Errorf("[account|HandleAddFriend] 更新 AccountFriend 发生错误, %s", err.Error())
		return
	}
	err = nil
	return
}

// PageApplyFriend 好友申请列表
func (d *Dao) PageApplyFriend(currId uint, param *model.AccountFriendPageParam) (resp []*model.Account, count int64, err error) {
	subQuery := config.AllConn.Db.Debug().Table("account_friend_apply as apply").
		Select("friend_id").
		Where("account_id = ? and status = 0", currId)
	db := config.AllConn.Db.Debug().Table("account").Where("account.id in (?)", subQuery)

	if err = db.Count(&count).Error; err != nil {
		logrus.Errorf("[account|PageApplyFriend] 发生错误, %s", err.Error())
		return
	}
	if err = db.Offset(int((param.Page - 1) * param.PageSize)).Limit(int(param.PageSize)).Find(&resp).Error; err != nil {
		logrus.Errorf("[account|PageApplyFriend] 分页查找account, err: [%+v]", err)
		return
	}
	return
}
