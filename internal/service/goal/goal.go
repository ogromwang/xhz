package goal

import (
	"xiaohuazhu/internal/dao/goal"
	"xiaohuazhu/internal/model"
	"xiaohuazhu/internal/util/result"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Service struct {
	goalDao *goal.Dao
}

func NewService() *Service {
	return &Service{
		goalDao: goal.New(),
	}
}

// Get 获取目标
func (s *Service) Get(ctx *gin.Context) {
	logrus.Infof("[goal|Set] 获取当前目标值")
	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	list, err := s.goalDao.Get(ctx, currUser.Id)
	if err != nil {
		result.ServerError(ctx)
		return
	}

	result.Ok(ctx, list)
}

// Set 设置目标
func (s *Service) Set(ctx *gin.Context) {
	logrus.Infof("[goal|Set] 修改目标值")
	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	var param = model.GoalSetDTO{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		result.Fail(ctx, "参数错误")
		return
	}

	set, err := s.goalDao.Set(ctx, &param, currUser.Id)
	if err != nil {
		result.ServerError(ctx)
		return
	}
	if !set {
		result.Fail(ctx, "更新失败")
		return
	}

	result.Success(ctx)
}
