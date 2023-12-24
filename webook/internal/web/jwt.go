package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type JwtHandler struct {
}

func (h *JwtHandler) setJwtToken(ctx *gin.Context, uid int64) {
	uc := UserClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserAgent: ctx.GetHeader("User-Agent"),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
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

var JwtKey = []byte("99c5468490C311Ee91Bb1A5958B90E3B")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
