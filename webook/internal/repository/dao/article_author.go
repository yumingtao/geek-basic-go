package dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleAuthorDao interface {
	Create(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, art Article) error
}

type ArticleGormAuthorDao struct {
	db *gorm.DB
}

func (a *ArticleGormAuthorDao) Create(ctx context.Context, art Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (a *ArticleGormAuthorDao) Update(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}

func NewArticleGormAuthorDao(db *gorm.DB) ArticleAuthorDao {
	return &ArticleGormAuthorDao{
		db: db,
	}
}
