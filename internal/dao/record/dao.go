package record

import (
	"xiaohuazhu/internal/config"
	"xiaohuazhu/internal/model"
)

type Dao struct {

}

func New() *Dao {
	return &Dao{}
}

func (d *Dao) Add(account *model.RecordMoney) (err error) {
	err = config.OrmDB.Save(account).Error
	return
}