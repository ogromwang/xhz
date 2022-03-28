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
	account.POST("", s.Account.Add)

	// 记录相关
	record := routerGroup.Group("record")
	record.POST("", s.Record.Push)
}