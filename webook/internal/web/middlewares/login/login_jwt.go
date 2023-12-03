package login

import (
	"encoding/gob"
	"geek-basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type JwtMiddlewareBuilder struct {
}

func (m *JwtMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
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
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			// 没有传token
			log.Println("没有传token")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			// 没有按照bearer token格式
			log.Println("bearer token格式不对")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := segs[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JwtKey, nil
		})
		if err != nil {
			// token不对
			log.Println("解析token报错")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			// token 解析出来，但是非法/已过期
			log.Println("token过期，请重新登录")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if uc.UserAgent != ctx.GetHeader("User-Agent") {
			// 后期监控告警要埋点，进入这个分支的大概率是攻击者
			log.Println("User-Agent不对")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expireTime := uc.ExpiresAt
		// 如果下边的代码成立，则 !token.Valid肯定是true，所以不用再判断，会被上边的代码拦截住
		/*if expireTime.Before(time.Now()) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}*/
		// 剩余过期时间小于50s就需要刷新
		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err := token.SignedString(web.JwtKey)
			ctx.Header("X-Jwt-Token", tokenStr)
			if err != nil {
				log.Println(err)
			}
		}
		// 登录成功之后，如果在context中设置好，后端不需要再去解析uc了
		ctx.Set("user", uc)
	}
}
