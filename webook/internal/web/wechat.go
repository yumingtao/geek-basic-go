package web

import (
	"fmt"
	"geek-basic-go/webook/internal/service"
	"geek-basic-go/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
)

type OAuth2WechatHandler struct {
	// 组合JwtHandler
	svc     wechat.Service
	userSvc service.UserService
	JwtHandler
	key             []byte
	stateCookieName string
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:             svc,
		userSvc:         userSvc,
		key:             JwtKey,
		stateCookieName: "jwt-state",
	}
}

func (o *OAuth2WechatHandler) OAuth2Url(ctx *gin.Context) {
	state := uuid.New()
	url, err := o.svc.AuthUrl(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "构造跳转URL失败",
			Code: 5,
			Data: nil,
		})
		return
	}
	err = o.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "构造跳转URL失败",
			Code: 5,
			Data: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	// 校验state
	err := o.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "非法请求",
			Code: 4,
		})
		return
	}
	// 校验或不校验都可以，如果值是空，后边调用微信接口会返回error
	code := ctx.Query("code")
	wechatInfo, err := o.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "授权码有误",
			Code: 4,
		})
		return
	}
	// 登录或注册逻辑（用户可能第一次登录）
	u, err := o.userSvc.FindOrCreatedByWechat(ctx, wechatInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	o.setJwtToken(ctx, u.Id)
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.OAuth2Url)
	g.Any("/callback", o.Callback)
}

func (o *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie(o.stateCookieName)
	if err != nil {
		return fmt.Errorf("无法获得cookie，%w", err)
	}
	var sc StateClaims
	_, err = jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return o.key, nil
	})
	if err != nil {
		return fmt.Errorf("解析token失败，%w", err)
	}

	if state != sc.state {
		return fmt.Errorf("state 不匹配")
	}
	return nil
}

func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	claims := StateClaims{
		state: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(o.key)
	//tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return err
	}
	// 这里直接set到了cookie，因为微信回来的时候是调到后端的回调接口
	ctx.SetCookie(o.stateCookieName, tokenString, 600,
		// 限制只在这个path生效
		"/oauth2/wechat/callback",
		// 同时要设置线上环境的域名，这里传“”
		// 这边由于是本地开发测试，把https禁止了，不过部署环境要开启https
		// httpOnly:true, 没有办法通过js来操作cookie
		"", false, true)
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	state string
}
