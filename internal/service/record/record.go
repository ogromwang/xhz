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

func (s *Service) Push(ctx *gin.Context) {
	logrus.Infof("[recordMoney | Push] 开始新建花销记录")
	var param = model.RecordMoneyDTO{}

	if err := ctx.ShouldBindJSON(&param); err != nil {
		result.Fail(ctx, "参数错误")
		return
	}

	var po = model.RecordMoney{
		AccountId: param.AccountId,
		Share:     param.Share,
		Money:     param.Money,
		Describe:  param.Describe,
		Image:     param.Image,
	}
	// 保存
	if err := s.recordDao.Add(&po); err != nil {
		result.Fail(ctx, err.Error())
		return
	}
	result.Success(ctx)
}