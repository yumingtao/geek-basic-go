package dao

import (
	"context"
	"github.com/ecodeclub/ekit/sqlx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var ErrSmsNotFound = gorm.ErrRecordNotFound

type SmsDao interface {
	Insert(ctx context.Context, sms Sms) error
	Get(ctx context.Context, interval time.Duration) (Sms, error)
	MarkSuccess(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64) error
}

const (
	// 初始放入重试数据库状态是0
	asyncStatusWaiting = iota
	// 达到重试阈值仍发送失败
	asyncStatusFailed
	// 经过重试发送成功
	asyncStatusSucceed
)

type GormSmsDao struct {
	db *gorm.DB
}

func NewGormSmsDao(db *gorm.DB) SmsDao {
	return &GormSmsDao{
		db: db,
	}
}

func (g *GormSmsDao) Insert(ctx context.Context, sms Sms) error {
	return g.db.WithContext(ctx).Create(&sms).Error
}

func (g *GormSmsDao) Get(ctx context.Context, interval time.Duration) (Sms, error) {
	// 取一条记录后要更新retryCnt，使用select for update对一条记录加锁
	var sms Sms
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		nowMilli := now.UnixMilli()
		endTime := now.Add(-interval).UnixMilli()
		//查找interval时间前最早的一条记录
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("utime < ? and status = ?", endTime, asyncStatusWaiting).First(&sms).Error
		if err != nil {
			return err
		}
		// 查询记录成功，更新retryCnt和时间
		err = tx.Model(&Sms{}).Where("id=?", sms.Id).
			Updates(map[string]any{
				"retry_cnt": gorm.Expr("retry_cnt + 1"),
				"utime":     nowMilli,
			}).Error
		return err
	})
	return sms, err
}

func (g *GormSmsDao) MarkSuccess(ctx context.Context, id int64) error {
	nowMilli := time.Now().UnixMilli()
	err := g.db.WithContext(ctx).Model(&Sms{}).Where("id=?", id).
		Updates(map[string]any{
			"status": asyncStatusSucceed,
			"utime":  nowMilli,
		}).Error
	return err
}

func (g *GormSmsDao) MarkFailed(ctx context.Context, id int64) error {
	nowMilli := time.Now().UnixMilli()
	err := g.db.WithContext(ctx).Model(&Sms{}).Where("id=? and `retry_cnt`>=`retry_max_cnt`", id).
		Updates(map[string]any{
			"status": asyncStatusFailed,
			"utime":  nowMilli,
		}).Error
	return err
}

type Sms struct {
	Id          int64
	SmsConfig   sqlx.JsonColumn[SmsConfig]
	RetryCnt    int
	RetryMaxCnt int
	status      uint8
	CTime       int64
	UTime       int64 `gorm:"index"`
}

type SmsConfig struct {
	Tpl     string
	Args    []string
	Numbers []string
}
