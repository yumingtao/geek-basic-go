package service

import (
	"context"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
}

type ArticleServiceImpl struct {
	repo repository.ArticleRepository
}

func (a *ArticleServiceImpl) Publish(ctx context.Context, art domain.Article) (int64, error) {
	panic(0)
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &ArticleServiceImpl{
		repo: repo,
	}
}

func (a *ArticleServiceImpl) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}
