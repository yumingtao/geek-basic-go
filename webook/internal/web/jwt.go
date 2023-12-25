package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type JwtHandler struct {
	signingMethod jwt.SigningMethod
	refreshKey    []byte
}

func newJwtHandler() JwtHandler {
	return JwtHandler{
		signingMethod: jwt.SigningMethodHS512,
		refreshKey:    []byte("99c5468490C311Ee91Bb1A5958B90E3A"),
	}
}

func (h *JwtHandler) setJwtToken(ctx *gin.Context, uid int64) {
	uc := UserClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserAgent: ctx.GetHeader("User-Agent"),
	}
	token := jwt.NewWithClaims(h.signingMethod, uc)
	signedString, err := token.SignedString(JwtKey)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.Header("X-Jwt-Token", signedString)
}

func (h *JwtHandler) setRefreshToken(ctx *gin.Context, uid int64) error {
	rc := RefreshClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	refreshToken := jwt.NewWithClaims(h.signingMethod, rc)
	signedString, err := refreshToken.SignedString(h.refreshKey)
	if err != nil {
		return err
	}
	ctx.Header("X-Refresh-Token", signedString)
	return nil
}

func ExtractToken(ctx *gin.Context) string {
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

var JwtKey = []byte("99c5468490C311Ee91Bb1A5958B90E3B")

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid int64
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
