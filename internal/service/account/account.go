package account

import (
	"xiaohuazhu/internal/dao/account"
	"xiaohuazhu/internal/model"
	"xiaohuazhu/internal/util/result"

	"github.com/gin-gonic/gin"
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
