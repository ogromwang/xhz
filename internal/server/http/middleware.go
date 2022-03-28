package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"xiaohuazhu/internal/model"
	"xiaohuazhu/internal/util/auth"
	"xiaohuazhu/internal/util/result"
)

func Auth() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader(model.X_TOKEN)
		if token == "" {
			logrus.Errorf("没有访问权限")
			result.Fail(ctx, "没有访问权限")
			ctx.Abort()
			return
		}
		dto, err := auth.ParseToken(token)
		if err != nil {
			logrus.Errorf("解析 token 失败, token: [%s], err: [%+v]", token, err)
			result.Fail(ctx, "没有访问权限")
			ctx.Abort()
			return
		}
		ctx.Set(model.CURR_USER, dto)
		ctx.Next()
	}
}
