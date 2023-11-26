package login

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MiddlewareBuilder struct {
}

func (m *MiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/login" {
			// 登录不需要校验
			println("登录不需要校验")
			return
		}
		method := ctx.Request.Method
		if path == "/users" && method == http.MethodPost {
			// 注册不需要校验
			println("注册不需要校验")
			return
		}
		sess := sessions.Default(ctx)
		if sess.Get("userId") == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
