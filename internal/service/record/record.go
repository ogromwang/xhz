package record

import (
	"github.com/gin-gonic/gin/binding"
	"os"
	"xiaohuazhu/internal/dao/record"
	"xiaohuazhu/internal/model"
	"xiaohuazhu/internal/util"
	"xiaohuazhu/internal/util/oss"
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

// Push push
func (s *Service) Push(ctx *gin.Context) {
	logrus.Infof("[recordMoney|Push] 开始新建记录")
	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	file, err := ctx.FormFile(model.FILE)
	if err != nil && err.Error() != "http: no such file" {
		logrus.Errorf("[record|Push] 读取上传文件发生错误, %s", err.Error())
		result.Fail(ctx, "上传文件失败")
		return
	}

	var path = ""
	// handle 文件
	if file != nil {
		// 8mb
		if file.Size > (8 << 20) {
			logrus.Errorf("[record|Push] 文件大小: %s", util.FormatFileSize(file.Size))
			result.Fail(ctx, "文件大小超过限制")
			return
		}

		var temp *os.File
		var compressTemp *os.File
		temp, err = util.NewImgTempPath(util.GetFileExt(file.Filename))
		compressTemp, err = util.NewImgTempPath(util.GetFileExt(file.Filename))
		if err != nil {
			logrus.Errorf("[record|Push] 创建临时目录异常")
			result.ServerError(ctx)
			return
		}

		// 处理上传的数据，写入临时
		defer os.Remove(temp.Name())
		defer temp.Close()
		if err = ctx.SaveUploadedFile(file, temp.Name()); err != nil {
			logrus.Errorf("[record|Push] 写入临时目录: [%s] 失败, %s", temp.Name(), err.Error())
			result.ServerError(ctx)
			return
		}
		// 压缩图片
		defer os.Remove(compressTemp.Name())
		defer compressTemp.Close()
		err = util.ImgFileResize(temp, compressTemp, 400)
		if err != nil {
			logrus.Errorf("[record|Push] 压缩图片时发生异常, err: [%s]", err.Error())
			result.ServerError(ctx)
			return
		}

		// 上传至 oss, 这里进行了压缩 io后，需要传递path重新读取？
		path, err = oss.PushObject(compressTemp.Name(), "picture")
		if err != nil {
			logrus.Errorf("[record|Push] OSS 上传头像失败: %s %s", temp.Name(), err.Error())
			result.Fail(ctx, "上传头像失败，请联系管理员")
			return
		}
	}

	// 新增记录，是否可见
	var param = model.RecordMoneyDTO{}
	if err = ctx.ShouldBindWith(&param, binding.FormMultipart); err != nil {
		result.Fail(ctx, "参数错误")
		return
	}
	var po = model.RecordMoney{
		AccountId: currUser.Id,
		Share:     param.Share,
		Money:     param.Money,
		Describe:  param.Describe,
		Image:     path,
	}
	// 保存
	if err = s.recordDao.Add(&po); err != nil {
		logrus.Errorf("[record|Push] DB 保存错误, %s", err.Error())
		result.ServerError(ctx)
		return
	}
	result.Success(ctx)
}

// RecordByFriends ...
func (s *Service) RecordByFriends(ctx *gin.Context) {
	logrus.Infof("[record|RecordByFriends] 查询记录")

	data := ctx.MustGet(model.CURR_USER)
	currUser := data.(*model.AccountDTO)

	var param = model.RecordPageParam{}
	if err := ctx.ShouldBindQuery(&param); err != nil {
		result.Fail(ctx, "参数错误")
		return
	}
	// 需要关联查询 record + account
	records, err := s.recordDao.RecordByFriends(&param, currUser)
	if err != nil {
		logrus.Errorf("[record|RecordByFriends] DB 查询错误, %s", err.Error())
		result.ServerError(ctx)
		return
	}

	var hasMore bool
	if len(records) > int(param.PageSize) {
		hasMore = true
		records = records[:param.PageSize]
	}

	for _, dto := range records {
		dto.ProfilePicture = oss.GetUrlByProtocol(dto.ProfilePicture)
		dto.Image = oss.GetUrlByProtocol(dto.Image)

	}
	result.OkWithMore(ctx, records, hasMore)
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
