package auth

import (
	"context"
	"geek-basic-go/webook/internal/service/sms"
	"github.com/golang-jwt/jwt/v5"
)

type SmsService struct {
	svc sms.Service
	key []byte
}

type SmsClaims struct {
	jwt.RegisteredClaims
	Tpl string
	// 可以额外加字段
}

func (s *SmsService) Send(ctx context.Context, tplToken string, args []string, numbers ...string) error {
	var claims SmsClaims
	_, err := jwt.ParseWithClaims(tplToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	return s.Send(ctx, claims.Tpl, args, numbers...)
}
