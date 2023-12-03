package login

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type MiddlewareBuilder struct {
}

func (m *MiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	// 注册一下time.Now
	gob.Register(time.Now())
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
		userId := sess.Get("userId")
		if userId == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		const updateTimeKey = "update_time"
		now := time.Now()
		// 拿出上一次的更新时间
		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Minute {
			// 第一次进来, 注意这里time.Now()需要在gob注册一下
			sess.Set(lastUpdateTime, now)
			// sess 设置是覆盖式，必须要把之前的userId一起设置，否则只有一个update time
			sess.Set("userId", userId)
			err := sess.Save()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
