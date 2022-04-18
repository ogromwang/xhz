package goal

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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

func (d *Dao) Get(ctx *gin.Context, currId uint) (list []*model.GoalGetDTO, err error) {
	// 获取当前 currId 的所有 目标
	poList := make([]*model.Goal, 0)
	list = make([]*model.GoalGetDTO, 0)

	// 小组
	if err = config.AllConn.Db.Table("goal").Where("? = any (account_ids)", currId).Find(&poList).Error; err != nil {
		logrus.Errorf("[goal|Set] 数据库错误, %s", err.Error())
		return
	}

	lock := sync.Mutex{}
	for _, goal := range poList {
		go func(g *model.Goal) {
			oneGoal, _ := d.getOneGoal(ctx, currId, g.ID)
			ptr := &model.GoalGetDTO{
				Id:         g.ID,
				CurrMoney:  g.Money,
				TotalMoney: oneGoal,
				Type:       g.Type,
			}

			lock.Lock()
			list = append(list, ptr)
			lock.Unlock()

		}(goal)
	}
	return
}

func (d *Dao) getOneGoal(_ *gin.Context, currId uint, goalId uint) (total float32, err error) {
	year := time.Now().Year()
	month := time.Now().Format("01")

	subQuery := config.AllConn.Db.Table("goal").
		Select("regexp_split_to_table(array_to_string(friend_ids, ','), ',')::int").
		Where("account_id = ?", currId).
		Where("id = ?", goalId)
	db := config.AllConn.Db.Table("record_money").
		Select("sum(money) as total").
		Where("account.id in (?)", subQuery).
		Where("date_part('year', created_at) = ?", year).
		Where("date_part('month', created_at) = ?", month).
		Group("date_part('month', created_at)")

	var res = struct {
		total float32
	}{}

	if err = db.Scan(&res).Error; err != nil {
		return
	}

	total = res.total
	return
}

func (d *Dao) Set(_ *gin.Context, g *model.GoalSetDTO, currId uint) (success bool, err error) {
	if g.Type == 1 {
		// 个人
		// 新增 | 修改
		create := model.Goal{
			AccountIds: pq.Int64Array([]int64{int64(currId)}),
			Leader:     currId,
			Money:      g.Money,
			Type:       g.Type,
		}
		if err = config.AllConn.Db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "leader"}, {Name: "type"}},
			DoUpdates: clause.AssignmentColumns([]string{"money"}),
		}).Create(&create).Error; err != nil {
			logrus.Errorf("[goal|Set] DB写入异常, %s", err.Error())
			return
		}

	} else if g.Type == 2 {
		if g.Id == 0 {
			return
		}

		var po = model.Goal{}
		// 小组
		// 先查询是否存在记录 | 是否是leader，不然不能修改
		if err = config.AllConn.Db.
			Where("leader = ? and type = 2", currId).
			Where("id = ?", g.Id).
			First(&po).
			Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				logrus.Errorf("[goal|Set] 没有找到消息或者没有权限, %s", err.Error())
				err = nil
				return
			}
			logrus.Errorf("[goal|Set] 数据库错误, %s", err.Error())
		}

		po.Money = g.Money
		if err = config.AllConn.Db.Save(&po).Error; err != nil {
			logrus.Errorf("[goal|Set] 数据库错误, %s", err.Error())
			return
		}

	}

	success = true
	return
}
