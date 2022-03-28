package account

import (
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
	if err = config.OrmDB.Find(&resp).Error; err != nil {
		logrus.Errorf("[account | list] 发生错误, %s", err.Error())
		return
	}
	return
}

func (d *Dao) Add(account *model.Account) (err error) {
	err = config.OrmDB.Save(account).Error
	return
}