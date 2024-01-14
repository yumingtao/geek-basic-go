package repository

import (
	"context"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository/cache"
	"geek-basic-go/webook/internal/repository/dao"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
}

type CachedArticleRepository struct {
	dao       dao.ArticleDao
	cache     cache.ArticleCache
	readerDao dao.ArticleReaderDao
	authorDao dao.ArticleAuthorDao
	db        *gorm.DB
}

func (c *CachedArticleRepository) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	// 先判断要不要查询缓存
	if offset == 0 && limit == 100 {
		res, err := c.cache.GetFirstPage(ctx, uid)
		if err == nil {
			return res, nil
		} else {
			// 记录日志
			// 缓存未命中，忽略
			// 网络问题
			// Redis 问题
		}
	}
	arts, err := c.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	res := slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		return c.toDomain(src)
	})
	go func() {
		// limit <= 100都可以缓存
		if offset == 0 && limit == 100 {
			// 缓存回写失败不一定是大问题，也可能是大问题
			err = c.cache.SetFirstPage(ctx, uid, res)
			if err != nil {
				// 网络抖动，记录日志，监控
			}
		}
	}()
	return res, nil
}

func (c *CachedArticleRepository) SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error {
	err := c.dao.SyncStatus(ctx, uid, id, status)
	if err == nil {
		err := c.cache.DeleteFirstPage(ctx, uid)
		if err != nil {
			// 记录日志
		}
	}
	return err
}

func NewArticleRepository(dao dao.ArticleDao, cache cache.ArticleCache) ArticleRepository {
	return &CachedArticleRepository{
		dao:   dao,
		cache: cache,
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
	id, err := c.dao.Sync(ctx, c.toEntity(art))
	if err == nil {
		err := c.cache.DeleteFirstPage(ctx, art.Author.Id)
		if err != nil {
			// 记录日志
		}
	}
	return id, err
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
	err := c.dao.UpdateById(ctx, c.toEntity(art))
	if err == nil {
		err := c.cache.DeleteFirstPage(ctx, art.Author.Id)
		if err != nil {
			// 记录日志
		}
	}
	return err
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	article := c.toEntity(art)
	id, err := c.dao.Insert(ctx, article)
	if err == nil {
		err := c.cache.DeleteFirstPage(ctx, art.Author.Id)
		if err != nil {
			// 记录日志
		}
	}
	return id, err
}

func (c *CachedArticleRepository) toEntity(art domain.Article) dao.Article {
	article := dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
		//Status:   uint8(art.Status),
	}
	return article
}

func (c *CachedArticleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.Id,
		},
		Status: domain.ArticleStatus(art.Status),
		Ctime:  time.UnixMilli(art.Ctime),
		Utime:  time.UnixMilli(art.Utime),
	}
}
