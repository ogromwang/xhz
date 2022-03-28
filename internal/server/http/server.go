package http

import (
	"github.com/gin-gonic/gin"
	"xiaohuazhu/internal/config"
	"xiaohuazhu/internal/service"
)

func New()  {
	r := gin.Default()
	s := service.New()
	router(r, s)

	r.Run(":" + config.AllConfig.Application.Port)
}