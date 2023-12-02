package web

import (
	"errors"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

// UserHandler
// 所有和用户相关路由都定义在这个handler上
// 用定义在UserHandler上的方法来作为路由的处理逻辑
type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
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
	ug.POST("/login", h.Login)
	ug.GET("/:id", h.Profile)
	ug.PUT("/:id", h.Edit)
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

func (h *UserHandler) LoginWithJwt(ctx *gin.Context) {
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
		uc := UserClaims{
			Uid: u.Id,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
		signedString, err := token.SignedString(JwtKey)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误！")
			return
		}
		ctx.Header("X-Jwt-Token", signedString)
		ctx.String(http.StatusOK, "恭喜，登录成功")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "登录失败："+err.Error())
	default:
		ctx.String(http.StatusOK, "系统错误！")
	}
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	//uc := ctx.MustGet("user").(UserClaims)
	ctx.String(http.StatusOK, "这是profile")
}

func (h *UserHandler) Edit(ctx *gin.Context) {

}

var JwtKey = []byte("99c5468490C311Ee91Bb1A5958B90E3B")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid int64
}
