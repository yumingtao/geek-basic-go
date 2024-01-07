package repository

import (
	"context"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository/dao"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
}

type CachedArticleRepository struct {
	dao dao.ArticleDao

	readerDao dao.ArticleReaderDao
	authorDao dao.ArticleAuthorDao

	db *gorm.DB
}

func NewArticleRepository(dao dao.ArticleDao) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

func NewArticleRepositoryV1(
	readerDao dao.ArticleReaderDao,
	authorDao dao.ArticleAuthorDao) *CachedArticleRepository {
	return &CachedArticleRepository{
		readerDao: readerDao,
		authorDao: authorDao,
	}
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	panic("Implement")
}

func (c *CachedArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// 防止后边业务panic
	defer tx.Rollback()
	authorDao := dao.NewArticleGormAuthorDao(tx)
	readerDao := dao.NewArticleGormReaderDao(tx)

	artEntity := c.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = authorDao.Update(ctx, artEntity)
	} else {
		id, err = authorDao.Create(ctx, artEntity)
	}

	if err != nil {
		return 0, err
	}
	artEntity.Id = id
	err = readerDao.UpsertV2(ctx, dao.PublishedArticle(artEntity))
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, nil
}

func (c *CachedArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	artEntity := c.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = c.authorDao.Update(ctx, artEntity)
	} else {
		id, err = c.authorDao.Create(ctx, artEntity)
	}

	if err != nil {
		return 0, err
	}
	artEntity.Id = id
	err = c.readerDao.Upsert(ctx, artEntity)
	return id, err
}

func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, c.toEntity(art))
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	article := c.toEntity(art)
	return c.dao.Insert(ctx, article)
}

func (c *CachedArticleRepository) toEntity(art domain.Article) dao.Article {
	article := dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
	return article
}
