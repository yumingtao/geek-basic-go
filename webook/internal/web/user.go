package web

import (
	"errors"
	"fmt"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/service"
	ijwt "geek-basic-go/webook/internal/web/jwt"
	"geek-basic-go/webook/pkg/logger"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
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
	errLogger       logger.ErrLogger
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, hdl ijwt.Handler, l logger.ErrLogger) *UserHandler {
	return &UserHandler{
		emailRexExp:     regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp:  regexp.MustCompile(passwordRegexPattern, regexp.None),
		birthDateRexExp: regexp.MustCompile(birthDateRegexPattern, regexp.None),
		svc:             svc,
		codeSvc:         codeSvc,
		Handler:         hdl,
		errLogger:       l,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	/*server.POST("/users", h.SignUp)
	server.POST("/users/login", h.Login)
	server.GET("/users/:id", h.Profile)
	server.PUT("/users/:id", h.Edit)*/
	// 分组注册路由
	ug := server.Group("/users")
	ug.POST("", h.SignUp)
	//ug.POST("/login", h.Login)
	ug.POST("/login", h.LoginWithJwt)
	ug.POST("/logout", h.LogoutWithJwt)
	ug.GET("/:id", h.Profile)
	ug.PUT("/:id", h.Edit)
	ug.PUT("/refresh_token", h.RefreshToken)
	// 短信验证码相关功能
	ug.POST("/login/sms/code", h.SendSmsLoginCode)
	ug.POST("/login/sms", h.VerifySmsCode)
}

func (h *UserHandler) SendSmsLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 校验手机号
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入正确手机号",
			Data: nil,
		})
		return
	}
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "短信发送成功",
		})
	case errors.Is(err, service.ErrCodeSentTooMany):
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "短信发送太频繁，请稍后再试",
		})
		// 埋点日志
		zap.L().Warn("频繁发送验证码")
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 打印日志
		log.Println(err)
	}
}

func (h *UserHandler) VerifySmsCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.errLogger.HandleError(err, "手机验证码验证失败")
		zap.L().Error("手机验证码验证失败:", zap.String("phone", req.Phone), zap.Error(err))
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码不正确，请重新输入",
		})
		// 补充日志
		return
	}
	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.errLogger.HandleError(err, "系统错误")
		return
	}
	err = h.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		h.errLogger.HandleError(err, "系统错误", logger.Field{
			Key: "origErr",
			Val: err,
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	// 内部类
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"ConfirmPassword"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	//isEmail, err := regexp.Match(emailRegexPattern, []byte(req.Email))
	/*if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}*/
	//isEmail := h.emailRexExp.Match([]byte(req.Email))
	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		h.errLogger.HandleError(err, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "邮箱格式不正确")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不一致")
		return
	}
	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		h.errLogger.HandleError(err, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含数字、特殊字符，并且长度不能小于8位")
		return
	}

	err = h.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	switch {
	case err == nil:
		ctx.String(http.StatusOK, "Hello, 恭喜注册成功")
	case errors.Is(err, service.ErrDuplicateEmail):
		ctx.String(http.StatusOK, "注册用户失败:"+err.Error())
	default:
		ctx.String(http.StatusOK, "系统错误！")
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
			h.errLogger.HandleError(err, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "恭喜，登录成功")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "登录失败："+err.Error())
	default:
		ctx.String(http.StatusOK, "系统错误！")
	}
}

func (h *UserHandler) LoginWithJwt(ctx *gin.Context) {

	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		h.errLogger.HandleError(err, "系统错误")
		return
	}
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
			ctx.String(http.StatusOK, "系统错误")
			h.errLogger.HandleError(err, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "恭喜，登录成功")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		h.errLogger.HandleError(err, "系统错误")
		ctx.String(http.StatusOK, "登录失败："+err.Error())
	default:
		ctx.String(http.StatusOK, "系统错误！")
	}
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	//uc := ctx.MustGet("user").(UserClaims)
	paramId := ctx.Param("id")
	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		ctx.String(http.StatusOK, "不是有效的用户id")
		h.errLogger.HandleError(err, "不是有效的用户id")
		return
	}
	u, err := h.svc.Profile(ctx, id)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, u)
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		NickName        string
		BirthDate       string
		PersonalProfile string
	}

	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if utf8.RuneCountInString(req.NickName) > nickNameMaxLen {
		ctx.String(http.StatusOK, "昵称允许最大长度"+fmt.Sprintf("%d", nickNameMaxLen)+", 请重新输入。")
		return
	}
	if utf8.RuneCountInString(req.PersonalProfile) > personalProfileMaxLen {
		ctx.String(http.StatusOK, "个人简介允许最大长度"+fmt.Sprintf("%d", personalProfileMaxLen)+", 请重新输入。")
		return
	}
	isBirthDate, err := h.birthDateRexExp.MatchString(req.BirthDate)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		h.errLogger.HandleError(err, "系统错误")
		return
	}
	if !isBirthDate {
		ctx.String(http.StatusOK, "不是有效生日，请重新输入。")
		return
	}

	paramId := ctx.Param("id")
	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		ctx.String(http.StatusOK, "不是有效的用户id")
		h.errLogger.HandleError(err, "不是有效的用户id")
		return
	}
	u, err := h.svc.Edit(ctx, domain.User{
		Id:              id,
		NickName:        req.NickName,
		BirthDate:       req.BirthDate,
		PersonalProfile: req.PersonalProfile,
	})

	switch {
	case err == nil:
		data := map[string]any{
			"id":              u.Id,
			"email":           u.Email,
			"nickName":        u.NickName,
			"birthDate":       u.BirthDate,
			"personalProfile": u.PersonalProfile,
		}
		ctx.JSON(http.StatusOK, data)
	default:
		ctx.String(http.StatusOK, "系统错误！")
	}
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
		h.errLogger.HandleError(err, "用户已登出")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.SetJwtToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		h.errLogger.HandleError(err, "设置Jwt token报错")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "刷新令牌成功",
	})
}

func (h *UserHandler) LogoutWithJwt(ctx *gin.Context) {
	err := h.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		h.errLogger.HandleError(err, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出登录成功",
	})
}
