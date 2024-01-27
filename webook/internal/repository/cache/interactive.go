package cache

import (
	"context"
	_ "embed"
	"fmt"
	"geek-basic-go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
)

const fieldReadCnt = "read_cnt"
const fieldLikeCnt = "like_cnt"
const fieldCollectCnt = "collect_cnt"

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, id int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error
}

type InteractiveRedisCache struct {
	client redis.Cmdable
}

func (i *InteractiveRedisCache) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	key := i.key(biz, bizId)
	res, err := i.client.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(res) == 0 {
		return domain.Interactive{}, ErrKeyNotExist
	}
	readCnt, _ := strconv.ParseInt(res[fieldReadCnt], 10, 64)
	likeCnt, _ := strconv.ParseInt(res[fieldLikeCnt], 10, 64)
	collectCnt, _ := strconv.ParseInt(res[fieldCollectCnt], 10, 64)
	return domain.Interactive{
		ReadCnt:    readCnt,
		LikeCnt:    likeCnt,
		CollectCnt: collectCnt,
	}, nil
}

func (i *InteractiveRedisCache) Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error {
	key := i.key(biz, bizId)
	err := i.client.HSet(ctx, key,
		fieldReadCnt, intr.ReadCnt,
		fieldLikeCnt, intr.LikeCnt,
		fieldCollectCnt, intr.CollectCnt,
	).Err()
	if err != nil {
		return err
	}
	return i.client.Expire(ctx, key, time.Minute*13).Err()
}

func NewInteractiveRedisCache(client redis.Cmdable) InteractiveCache {
	return &InteractiveRedisCache{
		client: client,
	}
}

func (i *InteractiveRedisCache) IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldCollectCnt, 1).Err()
}

func (i *InteractiveRedisCache) IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, 1).Err()
}

func (i *InteractiveRedisCache) DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, -1).Err()
}

func (i *InteractiveRedisCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := i.key(biz, bizId)
	// res 不需要特别处理，可以忽略
	//_, err := i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldReadCnt, 1).Int()
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldReadCnt, 1).Err()
}

func (i *InteractiveRedisCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
