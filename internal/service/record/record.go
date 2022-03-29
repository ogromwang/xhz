package record

import (
	"xiaohuazhu/internal/dao/record"
	"xiaohuazhu/internal/model"
	"xiaohuazhu/internal/util/result"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Service struct {
	recordDao *record.Dao
}

func NewService() *Service {
	return &Service{
		recordDao: record.New(),
	}
}

// Push
// todo-ogromwang : (2022.03.29 17:29) [待增强]
func (s *Service) Push(ctx *gin.Context) {
	logrus.Infof("[recordMoney|Push] 开始新建记录")
	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	// 新增记录，是否可见
	var param = model.RecordMoneyDTO{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		result.Fail(ctx, "参数错误")
		return
	}
	var po = model.RecordMoney{
		AccountId: currUser.Id,
		Share:     param.Share,
		Money:     param.Money,
		Describe:  param.Describe,
		Image:     param.Image,
	}
	// 保存
	if err := s.recordDao.Add(&po); err != nil {
		logrus.Errorf("[record|Push] DB 保存错误, %s", err.Error())
		result.ServerError(ctx)
		return
	}
	result.Success(ctx)
}

// RecordByFriends ...
func (s *Service) RecordByFriends(ctx *gin.Context) {
	logrus.Infof("[record|RecordByFriends] 查询记录")
	var param = model.RecordPageParam{}
	if err := ctx.ShouldBindQuery(&param); err != nil {
		result.Fail(ctx, "参数错误")
		return
	}
	// 需要关联查询 record + account
	records, err := s.recordDao.RecordByFriends(&param)
	if err != nil {
		logrus.Errorf("[record|RecordByFriends] DB 查询错误, %s", err.Error())
		result.ServerError(ctx)
		return
	}
	result.Ok(ctx, records)
}

// RecordByMe ...
func (s *Service) RecordByMe(ctx *gin.Context) {
	logrus.Infof("[record|RecordByMe] 查询个人 记录")
	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	var param = model.RecordPageParam{}
	if err := ctx.ShouldBindQuery(&param); err != nil {
		result.Fail(ctx, "参数错误")
		return
	}
	// 需要关联查询 record + account
	records, err := s.recordDao.RecordByMe(&param, currUser.Id)
	if err != nil {
		logrus.Errorf("[record|RecordByFriends] DB 查询错误, %s", err.Error())
		result.ServerError(ctx)
		return
	}
	result.Ok(ctx, records)
}
