package repository

import (
	"context"
	"geek-basic-go/webook/internal/domain"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, art domain.Article) error
}
