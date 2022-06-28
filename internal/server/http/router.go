package http

import (
	"github.com/gin-gonic/gin"

	"xiaohuazhu/internal/service"
)

func router(r *gin.Engine, s *service.Service) {
	// v1 api
	v1(r, s)
}

func v1(r *gin.Engine, s *service.Service) {
	// 使用网关控制
	// r.Use(Cors())

	// 8 Mib 这不能限制，是占用多少内存
	r.MaxMultipartMemory = 8 << 20
	routerGroup := r.Group("v1")

	// 绕过鉴权
	account := routerGroup.Group("account")
	account.POST("signup", s.Account.SignUp)
	account.POST("signin", s.Account.SignIn)

	// 用户相关
	account.Use(Auth())
	account.GET("profile", s.Account.Profile)
	account.GET("friends", s.Account.PageMyFriend)
	account.PUT("profile/picture", s.Account.ProfilePicture)
	account.GET("friends/find", s.Account.PageFindFriend)
	account.GET("friends/apply", s.Account.PageApplyFriend)
	account.POST("friends/apply", s.Account.ApplyAddFriend)
	account.POST("friends/handle", s.Account.HandleAddFriend)

	// 记录相关
	record := routerGroup.Group("record")
	record.Use(Auth())
	record.POST("", s.Record.Push)
	record.DELETE("", s.Record.Delete)
	record.GET("me", s.Record.RecordByMe)
	record.GET("all", s.Record.RecordByFriends)

	// 目标相关
	goal := routerGroup.Group("goal")
	goal.Use(Auth())
	goal.GET("", s.Goal.List)
	goal.POST("", s.Goal.Set)
}
