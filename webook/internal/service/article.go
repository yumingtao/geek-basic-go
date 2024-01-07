package service

import (
	"context"
	"errors"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository"
	"geek-basic-go/webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
}

type ArticleServiceImpl struct {
	repo       repository.ArticleRepository
	readerRepo repository.ArticleReaderRepository
	authorRepo repository.ArticleAuthorRepository
	l          logger.LoggerV1
}

func NewArticleServiceV1(
	readerRepo repository.ArticleReaderRepository,
	authorRepo repository.ArticleAuthorRepository,
	l logger.LoggerV1) *ArticleServiceImpl {
	return &ArticleServiceImpl{
		readerRepo: readerRepo,
		authorRepo: authorRepo,
		l:          l,
	}
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &ArticleServiceImpl{
		repo: repo,
	}
}

func (a *ArticleServiceImpl) Publish(ctx context.Context, art domain.Article) (int64, error) {
	a.repo.SyncV1(ctx, art)
	panic("Implement me")
}

func (a *ArticleServiceImpl) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	// 先操作制作库，再操作线上库
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = a.authorRepo.Update(ctx, art)
	} else {
		id, err = a.authorRepo.Create(ctx, art)
	}

	if err != nil {
		return 0, err
	}
	art.Id = id
	for i := 0; i < 3; i++ {
		// 线上库可能有数据也可能没有数据
		// insert or update
		err = a.readerRepo.Save(ctx, art)
		if err != nil {
			// 多接入一些tracing 工具
			a.l.Error("保存到制作库成功但是到线上库失败",
				logger.Int64("aid", art.Id),
				logger.Error(err))
		} else {
			return id, nil
		}
	}
	a.l.Error("保存到制作库成功但是到线上库失败, 重试次数耗尽",
		logger.Int64("aid", art.Id),
		logger.Error(err))
	return id, errors.New("保存到线上库失败, 重试次数耗尽")
}

func (a *ArticleServiceImpl) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}
