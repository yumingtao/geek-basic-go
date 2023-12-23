package repository

import (
	"context"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository/dao"
	"github.com/ecodeclub/ekit/sqlx"
	"time"
)

var ErrSmsNotFound = dao.ErrSmsNotFound

type SmsRepository interface {
	Add(ctx context.Context, sms domain.Sms) error
	GetEarliestSmsByInterval(ctx context.Context, interval time.Duration) (domain.Sms, error)
	MarkStatus(ctx context.Context, id int64, status bool) error
}

type AsyncSmsRepository struct {
	dao dao.SmsDao
}

func NewAsyncSmsRepository(dao dao.SmsDao) SmsRepository {
	return &AsyncSmsRepository{
		dao: dao,
	}
}

func (a *AsyncSmsRepository) Add(ctx context.Context, sms domain.Sms) error {
	err := a.dao.Insert(ctx, dao.Sms{
		RetryMaxCnt: sms.RetryMaxCnt,
		SmsConfig: sqlx.JsonColumn[dao.SmsConfig]{
			Val: dao.SmsConfig{
				Tpl:     sms.Tpl,
				Args:    sms.Args,
				Numbers: sms.Numbers,
			},
			Valid: true,
		},
	})
	return err
}

func (a *AsyncSmsRepository) GetEarliestSmsByInterval(ctx context.Context, interval time.Duration) (domain.Sms, error) {
	ds, err := a.dao.Get(ctx, interval)
	if err != nil {
		return domain.Sms{}, err
	}
	return domain.Sms{
		Id:          ds.Id,
		Tpl:         ds.SmsConfig.Val.Tpl,
		Args:        ds.SmsConfig.Val.Args,
		Numbers:     ds.SmsConfig.Val.Numbers,
		RetryMaxCnt: ds.RetryMaxCnt,
	}, nil
}

func (a *AsyncSmsRepository) MarkStatus(ctx context.Context, id int64, status bool) error {
	if status {
		return a.dao.MarkSuccess(ctx, id)
	}
	return a.dao.MarkFailed(ctx, id)
}
