package account

import (
	"xiaohuazhu/internal/dao/account"
	"xiaohuazhu/internal/model"
	"xiaohuazhu/internal/util/result"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Service struct {
	accountDao *account.Dao
}

func NewService() *Service {
	return &Service{
		accountDao: account.New(),
	}
}

func (s *Service) List(ctx *gin.Context) {
	list, err := s.accountDao.List()
	if err != nil {
		result.Fail(ctx, err.Error())
		return
	}
	var resp = make([]*model.AccountDTO, 0, len(list))
	var pr *model.AccountDTO
	for _, data := range list {
		pr = &model.AccountDTO{
			Id:       data.ID,
			Username: data.Username,
			CreateAt: data.CreatedAt,
		}
		resp = append(resp, pr)
	}

	result.Ok(ctx, resp)
}

func (s *Service) Add(ctx *gin.Context) {
	logrus.Infof("[account | Add] 开始创建用户")
	var param = model.AccountDTO{}

	if err := ctx.ShouldBindJSON(&param); err != nil {
		result.Fail(ctx, "参数错误")
		return
	}

	var po = model.Account{
		Username: param.Username,
		Password: param.Password,
	}
	// 保存
	if err := s.accountDao.Add(&po); err != nil {
		result.Fail(ctx, err.Error())
		return
	}
	result.Success(ctx)
}