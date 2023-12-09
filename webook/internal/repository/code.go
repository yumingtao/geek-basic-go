package repository

import (
	"context"
	"geek-basic-go/webook/internal/repository/cache"
)

var ErrCodeVerifyTooMany = cache.ErrVerifyCodeTooMany
var ErrorCodeSentTooMany = cache.ErrCodeSendTooMany

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: c,
	}
}

func (c *CodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}
