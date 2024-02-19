package login

import (
	"encoding/gob"
	ijwt "geek-basic-go/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

type JwtMiddlewareBuilder struct {
	ijwt.Handler
}

func NewJwtMiddlewareBuilder(hdl ijwt.Handler) *JwtMiddlewareBuilder {
	return &JwtMiddlewareBuilder{
		Handler: hdl,
	}
}

func (m *JwtMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	// 注册一下time.Now
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/login" || path == "/users/login/sms/code" || path == "/users/login/sms" ||
			path == "/oauth2/wechat/authurl" || path == "/oauth2/wechat/callback" {
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
		tokenStr := m.ExtractToken(ctx)
		var uc ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return ijwt.UcJwtKey, nil
		})
		if err != nil {
			// token不对
			log.Println("解析token报错", err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			// token 解析出来，但是非法/已过期
			// 是否可以在这里触发刷新token - 在这里刷新和自动刷新没有什么区别了
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

		// 因为使用refresh toke，这个部分不需要了
		//expireTime := uc.ExpiresAt

		// 如果下边的代码成立，则 !token.Valid肯定是true，所以不用再判断，会被上边的代码拦截住
		/*if expireTime.Before(time.Now()) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}*/
		// week-03 剩余过期时间小于50s就需要刷新
		// week-04 压测时，过期时间设置30分钟
		// 因为使用refresh toke，这个部分不需要了
		/*if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
			tokenStr, err := token.SignedString(web.UcJwtKey)
			ctx.Header("X-Jwt-Token", tokenStr)
			if err != nil {
				log.Println(err)
			}
		}*/
		// 登录成功之后，如果在context中设置好，后端不需要再去解析uc了

		// 这里查看下redis，用户是否登出
		err = m.CheckSession(ctx, uc.Ssid)
		if err != nil {
			// 用户已登出或者redis有问题
			log.Println("用户已登出")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 比较温和的做法，兼容redis异常，如果redis有问题，result会是默认的0值，这样允许用户登录继续使用系统
		/*if result > 0 {
			// 用户已登出
			log.Println("用户已登出")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}*/
		ctx.Set("user", uc)
	}
}
