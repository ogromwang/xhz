package http

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
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
			if validationError, ok := err.(*jwt.ValidationError); ok && validationError.Errors == jwt.ValidationErrorExpired {
				result.NoAuth(ctx, "会话已过期, 请重新登录")
				ctx.Abort()
				return
			}

			result.Fail(ctx, "没有访问权限")
			ctx.Abort()
			return
		}
		ctx.Set(model.CURR_USER, dto)
		ctx.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-TOKEN, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
