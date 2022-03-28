package http

import (
	"xiaohuazhu/internal/service"

	"github.com/gin-gonic/gin"
)

func router(r *gin.Engine, s *service.Service) {
	// v1 api
	v1(r, s)
}

func v1(r *gin.Engine, s *service.Service) {
	routerGroup := r.Group("v1")

	// 用户相关
	account := routerGroup.Group("account")
	account.GET("", s.Account.List)
	// 注册
	account.POST("signup", s.Account.SignUp)
	// 登录
	account.POST("signin", s.Account.SignIn)

	// 记录相关
	record := routerGroup.Group("record")
	record.POST("", s.Record.Push)
}
