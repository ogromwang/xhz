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
	account.POST("signup", s.Account.SignUp)
	account.POST("signin", s.Account.SignIn)
	// 下面的接口需要鉴权
	account.Use(Auth())
	account.GET("friends", s.Account.ListMyFriend)
	account.GET("friends/find", s.Account.PageFindFriend)
	account.GET("friends/apply", s.Account.ListApplyFriend)
	account.POST("friends/apply", s.Account.ApplyAddFriend)
	account.POST("friends/handle", s.Account.HandleAddFriend)

	// 记录相关
	record := routerGroup.Group("record")
	record.POST("", s.Record.Push)
}
