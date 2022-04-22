package goal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
	"time"
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

func (d *Dao) List(ctx *gin.Context, currId uint) (list []*model.GoalGetDTO, err error) {
	// 获取当前 currId 的所有 目标
	poList := make([]*model.Goal, 0)
	list = make([]*model.GoalGetDTO, 0)

	// 小组
	if err = config.AllConn.Db.Debug().Table("goal").Where("? = any (account_ids)", currId).Find(&poList).Error; err != nil {
		logrus.Errorf("[goal|Set] 数据库错误, %s", err.Error())
		return
	}

	group := sync.WaitGroup{}
	lock := sync.Mutex{}
	for _, goal := range poList {
		group.Add(1)
		go func(g *model.Goal) {
			oneGoal, _ := d.getOneGoal(ctx, currId, g.ID)
			ptr := &model.GoalGetDTO{
				Id:         g.ID,
				Name:       g.Name,
				AccountIds: []int64(g.AccountIds),
				Goal:       g.Money,
				CurrMoney:  oneGoal,
				Type:       g.Type,
			}

			lock.Lock()
			list = append(list, ptr)
			lock.Unlock()
			group.Done()

		}(goal)
	}

	group.Wait()
	return
}

func (d *Dao) getOneGoal(_ *gin.Context, currId uint, goalId uint) (total float32, err error) {
	year := time.Now().Year()
	month := time.Now().Format("01")

	subQuery := config.AllConn.Db.Table("goal").
		Select("regexp_split_to_table(array_to_string(account_ids, ','), ',')::int").
		Where("account_id = ?", currId).
		Where("id = ?", goalId)
	db := config.AllConn.Db.Table("record_money").
		Select("sum(money) as total").
		Where("account_id in (?)", subQuery).
		Where("date_part('year', created_at) = ?", year).
		Where("date_part('month', created_at) = ?", month).
		Group("date_part('month', created_at)")

	if err = db.Debug().Scan(&total).Error; err != nil {
		return
	}
	logrus.Infof("值为: %v", total)
	return
}

func (d *Dao) findGoal(leaderId uint, typ int, goalId uint) ([]*model.Goal, error) {
	var poList = make([]*model.Goal, 0)
	// 小组
	// 先查询是否存在记录 | 是否是leader，不然不能修改
	db := config.AllConn.Db
	if goalId != 0 {
		db = db.Where("id = ?", goalId)
	} else {
		db = db.Where("leader = ? and type = ?", leaderId, typ)
	}

	if err := db.Find(&poList).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.Errorf("[goal|findGoal] 没有找到消息或者没有权限, %s", err.Error())
			return nil, nil
		}
		logrus.Errorf("[goal|findGoal] 数据库错误, %s", err.Error())
		return nil, err
	}
	return poList, nil
}

func (d *Dao) Set(_ *gin.Context, g *model.GoalSetDTO, currId uint) (success bool, err error) {
	var list = make([]*model.Goal, 0)
	// 通过主键直接拉
	if list, err = d.findGoal(currId, g.Type, g.Id); err != nil {
		logrus.Errorf("[goal|Set] DB写入异常, %s", err.Error())
		return
	}
	if len(list) != 1 {
		err = fmt.Errorf("没有发现记录")
		return
	}

	if err = d.saveUpdate(list[0].ID, currId, &model.Goal{
		Money:      g.Money,
		Type:       g.Type,
		AccountIds: g.AccountIds,
	}, config.AllConn.Db); err != nil {
		logrus.Errorf("[goal|Set] DB写入异常, %s", err.Error())
		return
	}

	success = true
	return
}

func (d *Dao) saveUpdate(id uint, currId uint, g *model.Goal, db *gorm.DB) error {
	var err error

	// 新增、修改
	create := model.Goal{
		Model: gorm.Model{
			ID: id,
		},
		AccountIds: g.AccountIds,
		Leader:     currId,
		Money:      g.Money,
		Type:       g.Type,
		Name:       g.Name,
	}
	if err = db.Debug().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"money", "name"}),
	}).Create(&create).Error; err != nil {
		logrus.Errorf("[goal|Set] DB写入异常, %s", err.Error())
		return err
	}
	return nil
}

func (d *Dao) Create(_ *gin.Context, g *model.GoalCreateDTO, currId uint, db *gorm.DB) (success bool, err error) {
	// type = 1 只能有一条记录
	// type = 2 一个人只能加入2个目标

	goalList := make([]*model.Goal, 0)
	goalList, err = d.findGoal(currId, g.Type, 0)

	if g.Type == 1 && len(goalList) != 0 {
		// 不能创建
		return
	}

	if g.Type == 2 && len(goalList) >= 2 {
		// 不能创建
		return
	}

	if err = d.saveUpdate(0, currId, &model.Goal{
		AccountIds: []int64{int64(currId)},
		Leader:     currId,
		Money:      g.Money,
		Type:       g.Type,
		Name:       g.Name,
	}, db); err != nil {
		logrus.Errorf("[goal|Set] DB写入异常, %s", err.Error())
		return
	}

	success = true
	return
}
