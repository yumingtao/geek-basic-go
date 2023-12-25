package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"strings"
	"time"
)

type RedisJwtHandler struct {
	signingMethod jwt.SigningMethod
	client        redis.Cmdable
	rcExpiration  time.Duration
}

func NewRedisJwtHandler(client redis.Cmdable) Handler {
	return &RedisJwtHandler{
		client:        client,
		signingMethod: jwt.SigningMethodHS512,
		rcExpiration:  time.Minute * 24 * 7,
	}
}

func (h *RedisJwtHandler) CheckSession(ctx *gin.Context, ssid string) error {
	result, err := h.client.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	if err != nil {
		return err
	}

	if result > 0 {
		// 用户已登出
		log.Println("用户已登出")
		return errors.New("token 无效")
	}
	return nil
}

func (h *RedisJwtHandler) ExtractToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		// 没有传token
		log.Println("没有传token")
		return ""
	}

	segs := strings.Split(authCode, " ")
	if len(segs) != 2 {
		// 没有按照bearer token格式
		log.Println("bearer token格式不对")
		return ""
	}

	return segs[1]
}

func (h *RedisJwtHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("X-Jwt-Token", "")
	ctx.Header("X-Refresh-Token", "")
	uc := ctx.MustGet("user").(UserClaims)
	return h.client.Set(ctx, fmt.Sprintf("users:ssid:%s", uc.Ssid), "", h.rcExpiration).Err()
}

func (h *RedisJwtHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.setRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	return h.SetJwtToken(ctx, uid, ssid)
}

func (h *RedisJwtHandler) SetJwtToken(ctx *gin.Context, uid int64, ssid string) error {
	uc := UserClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserAgent: ctx.GetHeader("User-Agent"),
	}
	token := jwt.NewWithClaims(h.signingMethod, uc)
	signedString, err := token.SignedString(UcJwtKey)
	if err != nil {
		return err
	}
	ctx.Header("X-Jwt-Token", signedString)
	return nil
}

func (h *RedisJwtHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	rc := RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.rcExpiration)),
		},
	}
	refreshToken := jwt.NewWithClaims(h.signingMethod, rc)
	signedString, err := refreshToken.SignedString(RcJwtKey)
	if err != nil {
		return err
	}
	ctx.Header("X-Refresh-Token", signedString)
	return nil
}

var UcJwtKey = []byte("99c5468490C311Ee91Bb1A5958B90E3B")
var RcJwtKey = []byte("99c5468490C311Ee91Bb1A5958B90E3A")

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	Ssid      string
	UserAgent string
}
