package dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleReaderDao interface {
	Upsert(ctx context.Context, art Article) error
	UpsertV2(ctx context.Context, art PublishedArticle) error
}

type ArticleGormReaderDao struct {
	db *gorm.DB
}

func NewArticleGormReaderDao(db *gorm.DB) ArticleReaderDao {
	return &ArticleGormReaderDao{
		db: db,
	}
}

func (a *ArticleGormReaderDao) UpsertV2(ctx context.Context, art PublishedArticle) error {
	//TODO implement me
	panic("implement me")
}

func (a *ArticleGormReaderDao) Upsert(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}
