package service

import (
	"context"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository"
)

type TopLikedService interface {
	GetTopLiked(ctx context.Context) ([]domain.Article, error)
}

type TopLikedServiceImpl struct {
	repo repository.TopLikedRepository
	n    int
}

func NewTopLikedServiceImpl(repo repository.TopLikedRepository) TopLikedService {
	return &TopLikedServiceImpl{
		repo: repo,
		n:    10,
	}
}

func (t *TopLikedServiceImpl) GetTopLiked(ctx context.Context) ([]domain.Article, error) {
	return t.repo.GetTopN(ctx, t.n)
}
