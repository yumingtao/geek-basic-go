package web

import (
	"errors"
	"fmt"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/service"
	ijwt "geek-basic-go/webook/internal/web/jwt"
	"geek-basic-go/webook/pkg/ginx"
	"geek-basic-go/webook/pkg/logger"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"unicode/utf8"
)

const (
	emailRegexPattern     = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern  = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	birthDateRegexPattern = `^(19|20)\d{2}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`
	nickNameMaxLen        = 20
	personalProfileMaxLen = 150
	bizLogin              = "login"
)

// UserHandler
// 所有和用户相关路由都定义在这个handler上
// 用定义在UserHandler上的方法来作为路由的处理逻辑
type UserHandler struct {
	// 使用组合，JwtHandler
	ijwt.Handler
	emailRexExp     *regexp.Regexp
	passwordRexExp  *regexp.Regexp
	birthDateRexExp *regexp.Regexp
	svc             service.UserService
	codeSvc         service.CodeService
	l               logger.LoggerV1
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, hdl ijwt.Handler, l logger.LoggerV1) *UserHandler {
	return &UserHandler{
		emailRexExp:     regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp:  regexp.MustCompile(passwordRegexPattern, regexp.None),
		birthDateRexExp: regexp.MustCompile(birthDateRegexPattern, regexp.None),
		svc:             svc,
		codeSvc:         codeSvc,
		Handler:         hdl,
		l:               l,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	/*server.POST("/users", h.SignUp)
	server.POST("/users/login", h.Login)
	server.GET("/users/:id", h.Profile)
	server.PUT("/users/:id", h.Edit)*/
	// 分组注册路由
	ug := server.Group("/users")
	ug.POST("", ginx.WrapBody[SignUpReq](h.SignUp))
	//ug.POST("/login", h.Login)
	ug.POST("/login", ginx.WrapBody(h.LoginWithJwt))
	ug.POST("/logout", h.LogoutWithJwt)
	ug.GET("/profile", ginx.WrapClaims(h.Profile))
	ug.PUT("/edit", ginx.WrapBodyAndClaims(h.Edit))
	ug.PUT("/refresh_token", h.RefreshToken)
	// 短信验证码相关功能
	ug.POST("/login/sms/code", ginx.WrapBody(h.SendSmsLoginCode))
	ug.POST("/login/sms", ginx.WrapBody(h.VerifySmsCode))
}

func (h *UserHandler) SendSmsLoginCode(ctx *gin.Context, req SendSmsCodeReq) (ginx.Result, error) {
	// 校验手机号
	if req.Phone == "" {
		return ginx.Result{
			Code: 4,
			Msg:  "请输入正确手机号",
			Data: nil,
		}, nil
	}
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch {
	case err == nil:
		return ginx.Result{
			Msg: "短信发送成功",
		}, nil
	case errors.Is(err, service.ErrCodeSentTooMany):
		// 埋点日志
		zap.L().Warn("频繁发送验证码")
		return ginx.Result{
			Code: 4,
			Msg:  "短信发送太频繁，请稍后再试",
		}, nil
	default:
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
}

func (h *UserHandler) VerifySmsCode(ctx *gin.Context, req VerifySmsCodeReq) (ginx.Result, error) {
	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		//zap.L().Error("手机验证码验证失败:", zap.String("phone", req.Phone), zap.Error(err))
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	if !ok {
		return ginx.Result{
			Code: 4,
			Msg:  "验证码不正确，请重新输入",
		}, nil
	}
	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	err = h.SetLoginToken(ctx, u.Id)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		Msg: "登录成功",
	}, nil
}

func (h *UserHandler) SignUp(ctx *gin.Context, req SignUpReq) (ginx.Result, error) {
	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	if !isEmail {
		return ginx.Result{
			Code: 4,
			Msg:  "邮箱格式不正确",
		}, err
	}

	if req.Password != req.ConfirmPassword {
		return ginx.Result{
			Code: 4,
			Msg:  "两次输入密码不一致",
		}, err
	}
	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	if !isPassword {
		return ginx.Result{
			Code: 4,
			Msg:  "密码必须包含数字、特殊字符，并且长度不能小于8位",
		}, err
	}

	err = h.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	switch {
	case err == nil:
		ctx.String(http.StatusOK, "Hello, 恭喜注册成功")
		return ginx.Result{
			Msg: "Hello, 恭喜注册成功",
		}, nil
	case errors.Is(err, service.ErrDuplicateEmail):
		return ginx.Result{
			Msg: "注册用户失败:" + err.Error(),
		}, nil
	default:
		return ginx.Result{
			Msg: "系统错误！",
		}, err
	}
}

// session logout
/*func (h *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	err := sess.Save()
	if err != nil {
		return
	}
}*/

func (h *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch {
	case err == nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			MaxAge:   900, //15分钟
			HttpOnly: true,
			Secure:   true,
		})
		err := sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误！")
			return
		}
		ctx.String(http.StatusOK, "恭喜，登录成功")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "登录失败："+err.Error())
	default:
		ctx.String(http.StatusOK, "系统错误！")
	}
}

func (h *UserHandler) LoginWithJwt(ctx *gin.Context, req LoginReq) (ginx.Result, error) {
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	userAgent := ctx.GetHeader("User-Agent")
	log.Println("User-Agent:", userAgent)
	switch {
	case err == nil:
		/*uc := UserClaims{
			Uid: u.Id,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			},
			UserAgent: userAgent,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
		signedString, err := token.SignedString(UcJwtKey)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误！")
			return
		}
		ctx.Header("X-Jwt-Token", signedString)*/
		err := h.SetLoginToken(ctx, u.Id)
		if err != nil {
			return ginx.Result{
				Code: 5,
				Msg:  "系统错误",
			}, err
		}
		return ginx.Result{
			Msg: "恭喜，登录成功",
		}, nil
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		return ginx.Result{
			Msg: "登录失败" + err.Error(),
		}, nil
	default:
		return ginx.Result{
			Msg: "系统错误！",
		}, nil
	}
}

func (h *UserHandler) Profile(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	//uc := ctx.MustGet("user").(UserClaims)
	/*paramId := ctx.Param("id")
	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		ctx.String(http.StatusOK, "不是有效的用户id")
		h.l.HandleError(err, "不是有效的用户id")
		return
	}*/
	u, err := h.svc.Profile(ctx, uc.Uid)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误：" + err.Error(),
		}, err
	}
	return ginx.Result{
		Data: u,
	}, nil
}

func (h *UserHandler) Edit(ctx *gin.Context, req EditReq, uc ijwt.UserClaims) (ginx.Result, error) {
	if utf8.RuneCountInString(req.NickName) > nickNameMaxLen {
		return ginx.Result{
			Code: 4,
			Msg:  "昵称允许最大长度" + fmt.Sprintf("%d", nickNameMaxLen) + ", 请重新输入。",
		}, nil
	}
	if utf8.RuneCountInString(req.PersonalProfile) > personalProfileMaxLen {
		return ginx.Result{
			Code: 4,
			Msg:  "个人简介允许最大长度" + fmt.Sprintf("%d", personalProfileMaxLen) + ", 请重新输入。",
		}, nil
	}
	isBirthDate, err := h.birthDateRexExp.MatchString(req.BirthDate)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	if !isBirthDate {
		return ginx.Result{
			Code: 4,
			Msg:  "不是有效生日，请重新输入。",
		}, nil
	}

	/*paramId := ctx.Param("id")
	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		ctx.String(http.StatusOK, "不是有效的用户id")
		h.l.HandleError(err, "不是有效的用户id")
		return
	}*/
	u, err := h.svc.Edit(ctx, domain.User{
		Id:              uc.Uid,
		NickName:        req.NickName,
		BirthDate:       req.BirthDate,
		PersonalProfile: req.PersonalProfile,
	})

	if err == nil {
		data := map[string]any{
			"id":              u.Id,
			"email":           u.Email,
			"nickName":        u.NickName,
			"birthDate":       u.BirthDate,
			"personalProfile": u.PersonalProfile,
		}
		return ginx.Result{
			Data: data,
			Msg:  "编辑用户信息成功。",
		}, nil
	}
	return ginx.Result{
		Code: 5,
		Msg:  "系统错误！",
	}, nil
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	// 前端在Authorization中带上refresh token
	tokenStr := h.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RcJwtKey, nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// 用户已登出或者redis有问题
		log.Println("用户已登出")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.SetJwtToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "刷新令牌成功",
	})
}

func (h *UserHandler) LogoutWithJwt(ctx *gin.Context) {
	err := h.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "退出登录成功",
	})
}
