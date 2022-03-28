package account

import (
	"github.com/sirupsen/logrus"
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

// ListMyFriend 我的好友
func (s *Service) ListMyFriend(ctx *gin.Context) {
	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	list, err := s.accountDao.ListFriend(currUser.Id)
	if err != nil {
		logrus.Errorf("[account|ListFriend] 发生错误, %s", err.Error())
		result.ServerError(ctx)
		return
	}

	result.Ok(ctx, s.transDTO(&list))
}

// PageFindFriend 查找用户，分页
func (s *Service) PageFindFriend(ctx *gin.Context) {
	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	logrus.Infof("[account|PageFindFriend] 寻找账号")
	var param = model.AccountFriendPageParam{}
	if err := ctx.ShouldBindQuery(&param); err != nil {
		result.Fail(ctx, "参数错误")
		return
	}

	accounts, err := s.accountDao.PageAccount(currUser.Id, &param)
	if err != nil {
		logrus.Errorf("[account|PageFindFriend] 发生错误, %s", err.Error())
		result.ServerError(ctx)
		return
	}
	result.Ok(ctx, s.transDTO(&accounts))
}

// ApplyAddFriend 申请添加好友
func (s *Service) ApplyAddFriend(ctx *gin.Context) {
	// 1. body 传递待添加参数
	logrus.Infof("[account|ApplyAddFriend] 申请添加好友")

	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	var param = model.ApplyAddFriendParam{}
	if err := ctx.ShouldBindQuery(&param); err != nil {
		result.Fail(ctx, "参数错误")
		return
	}
	// 2. 判断用户是否存在
	if _, err := s.accountDao.GetByUsernameOrId("", param.Id, true); err != nil {
		logrus.Errorf("[account|ApplyAddFriend] 发生错误, %s", err.Error())
		result.Fail(ctx, "用户不存在")
		return
	}

	// 3. 写入 account_friend_apply 中，并通知对方 todo，这个先不做，刷新就行
	if err := s.accountDao.ApplyAddFriend(param.Id, currUser.Id); err != nil {
		logrus.Errorf("[account|ApplyAddFriend] 发生错误, %s", err.Error())
		result.ServerError(ctx)
		return
	}
	result.Success(ctx)
}

// transDTO 转换为 DTO 返回
func (s *Service) transDTO(accounts *[]*model.Account) []*model.AccountDTO {
	var resp = make([]*model.AccountDTO, 0, len(*accounts))
	var pr *model.AccountDTO
	for _, data := range *accounts {
		pr = &model.AccountDTO{
			Id:       data.ID,
			Username: data.Username,
			CreateAt: data.CreatedAt,
		}
		resp = append(resp, pr)
	}
	return resp
}
