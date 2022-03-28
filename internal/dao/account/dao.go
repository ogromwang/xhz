package account

import (
	"gorm.io/gorm"
	"xiaohuazhu/internal/config"
	"xiaohuazhu/internal/model"

	"github.com/sirupsen/logrus"
)

type Dao struct {
}

func New() *Dao {
	return &Dao{}
}

func (d *Dao) List() (resp []*model.Account, err error) {
	if err = config.AllConn.Db.Find(&resp).Error; err != nil {
		logrus.Errorf("[account | list] 发生错误, %s", err.Error())
		return
	}
	return
}

func (d *Dao) Add(account *model.Account) (err error) {
	err = config.AllConn.Db.Save(account).Error
	return
}

func (d *Dao) GetByUsername(username string, onlyExist bool) (po *model.Account, err error) {
	po = &model.Account{}
	db := config.AllConn.Db.Where("username = ?", username)
	if onlyExist {
		db.Select("id")
	}
	err = db.Limit(1).First(po).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return po, nil
}
