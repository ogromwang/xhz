package account

import (
	"os"
	"strings"
	"xiaohuazhu/internal/dao/goal"
	"xiaohuazhu/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"xiaohuazhu/internal/dao/account"
	"xiaohuazhu/internal/model"
	"xiaohuazhu/internal/util/oss"
	"xiaohuazhu/internal/util/result"
)

type Service struct {
	accountDao *account.Dao
	goalDao    *goal.Dao
}

func NewService() *Service {
	return &Service{
		accountDao: account.New(),
		goalDao:    goal.New(),
	}
}

// PageMyFriend 我的好友
func (s *Service) PageMyFriend(ctx *gin.Context) {
	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	logrus.Infof("[account|PageMyFriend] 寻找账号")
	var param = model.AccountFriendPageParam{}
	if err := ctx.ShouldBindQuery(&param); err != nil {
		result.Fail(ctx, "参数错误")
		return
	}

	list, total, err := s.accountDao.PageFriend(currUser.Id, &param)
	if err != nil {
		logrus.Errorf("[account|PageFriend] DB 发生错误, %s", err.Error())
		result.ServerError(ctx)
		return
	}

	result.OkWithTotal(ctx, s.transDTO(&list), total)
}

func (s *Service) Profile(ctx *gin.Context) {
	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	list, err := s.accountDao.List([]int64{int64(currUser.Id)})
	if err != nil {
		logrus.Errorf("[account|Profile] DB 发生错误, %s", err.Error())
		result.ServerError(ctx)
		return
	}
	if len(list) != 1 {
		logrus.Warnf("[account|Profile] 警告, 数据异常")
		result.Fail(ctx, "没有该用户")
		return
	}
	user := list[0]
	result.Ok(ctx, model.AccountDTO{
		Id:             user.ID,
		Username:       user.Username,
		ProfilePicture: oss.GetUrlByProtocol(user.ProfilePicture),
		CreateAt:       user.CreatedAt,
	})
}

// ProfilePicture PUT 修改
func (s *Service) ProfilePicture(ctx *gin.Context) {
	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	file, err := ctx.FormFile(model.FILE)
	if err != nil {
		logrus.Errorf("[account|ProfilePicture] 读取上传文件发生错误, %s", err.Error())
		result.Fail(ctx, "上传文件失败")
		return
	}

	// 3mb
	if file.Size > (3 << 20) {
		logrus.Errorf("[account|ProfilePicture] 文件大小: %s", util.FormatFileSize(file.Size))
		result.Fail(ctx, "文件大小超过限制")
		return
	}
	var temp *os.File
	var compressTemp *os.File
	temp, err = util.NewImgTempPath(util.GetFileExt(file.Filename))
	compressTemp, err = util.NewImgTempPath(util.GetFileExt(file.Filename))
	if err != nil {
		logrus.Errorf("[account|ProfilePicture] 创建临时目录异常")
		result.ServerError(ctx)
		return
	}

	// 处理上传的数据，写入临时
	defer os.Remove(temp.Name())
	defer temp.Close()
	if err = ctx.SaveUploadedFile(file, temp.Name()); err != nil {
		logrus.Errorf("[account|ProfilePicture] 写入临时目录: [%s] 失败, %s", temp.Name(), err.Error())
		result.ServerError(ctx)
		return
	}
	// 压缩图片
	defer os.Remove(compressTemp.Name())
	defer compressTemp.Close()
	err = util.ImgFileResize(temp, compressTemp, 200)
	if err != nil {
		logrus.Errorf("[account|ProfilePicture] 压缩图片时发生异常, err: [%s]", err.Error())
		result.ServerError(ctx)
		return
	}

	// 上传至 oss, 这里进行了压缩 io后，需要传递path重新读取？
	path, err := oss.PushObject(compressTemp.Name(), "picture")
	if err != nil {
		logrus.Errorf("[account|ProfilePicture] OSS 上传头像失败: %s %s", temp.Name(), err.Error())
		result.Fail(ctx, "上传头像失败，请联系管理员")
		return
	}
	// 回写 DB
	if err = s.accountDao.UpdatePicture(currUser.Id, path); err != nil {
		logrus.Errorf("[account|ProfilePicture] 更新 DB 头像失败: %s %s", temp.Name(), err.Error())
		result.Fail(ctx, "更新信息失败")
		return
	}
	result.Success(ctx)
}

// PageFindFriend 查找用户，分页，需要过滤掉已经有的
func (s *Service) PageFindFriend(ctx *gin.Context) {
	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	logrus.Infof("[account|PageFindFriend] 寻找账号")
	var param = model.AccountFriendPageParam{}
	if err := ctx.ShouldBindQuery(&param); err != nil {
		result.Fail(ctx, "参数错误")
		return
	}
	if strings.Trim(param.Username, " ") == "" {
		result.Fail(ctx, "请填写名称")
		return
	}

	accounts, total, err := s.accountDao.PageFindAccount(currUser.Id, &param)
	if err != nil {
		logrus.Errorf("[account|PageFindFriend] 发生错误, %s", err.Error())
		result.ServerError(ctx)
		return
	}

	result.OkWithTotal(ctx, s.transDTO(&accounts), total)
}

// ApplyAddFriend 申请添加好友
func (s *Service) ApplyAddFriend(ctx *gin.Context) {
	// 1. body 传递待添加参数
	logrus.Infof("[account|ApplyAddFriend] 申请添加好友")

	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	var param = model.ApplyAddFriendParam{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
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

// HandleAddFriend 处理申请添加好友
func (s *Service) HandleAddFriend(ctx *gin.Context) {
	// 1. body 传递待添加参数
	logrus.Infof("[account|HandleAddFriend] 处理申请添加好友")

	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	var param = model.HandleAddFriendParam{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		result.Fail(ctx, "参数错误")
		return
	}
	// 2. 判断用户是否存在
	if _, err := s.accountDao.GetByUsernameOrId("", param.Id, true); err != nil {
		logrus.Errorf("[account|ApplyAddFriend] 发生错误, %s", err.Error())
		result.Fail(ctx, "用户不存在")
		return
	}

	if err := s.accountDao.HandleAddFriend(param.Id, currUser.Id, param.Status); err != nil {
		logrus.Errorf("[account|HandleAddFriend] DB 发生错误, %s", err.Error())
		result.ServerError(ctx)
		return
	}
	result.Success(ctx)
}

// PageApplyFriend 待处理的申请
func (s *Service) PageApplyFriend(ctx *gin.Context) {
	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)
	logrus.Infof("[account|PageApplyFriend] 显示待处理的申请")

	var param = model.AccountFriendPageParam{}
	if err := ctx.ShouldBindQuery(&param); err != nil {
		result.Fail(ctx, "参数错误")
		return
	}

	// 通过自己的id，查询 apple 表中的数据
	friend, total, err := s.accountDao.PageApplyFriend(currUser.Id, &param)
	if err != nil {
		logrus.Errorf("[account|PageApplyFriend] 发生错误, %s", err.Error())
		result.ServerError(ctx)
		return
	}
	result.OkWithTotal(ctx, s.transDTO(&friend), total)
}

// transDTO 转换为 DTO 返回
func (s *Service) transDTO(accounts *[]*model.Account) []*model.AccountDTO {
	var resp = make([]*model.AccountDTO, 0, len(*accounts))
	var pr *model.AccountDTO
	for _, data := range *accounts {
		pr = &model.AccountDTO{
			Id:             data.ID,
			ProfilePicture: oss.GetUrlByProtocol(data.ProfilePicture),
			Username:       data.Username,
			CreateAt:       data.CreatedAt,
		}
		resp = append(resp, pr)
	}
	return resp
}
