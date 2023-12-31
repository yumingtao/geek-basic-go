package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
}

type ArticleGormDao struct {
	db *gorm.DB
}

func (a *ArticleGormDao) Sync(ctx context.Context, art Article) (int64, error) {
	var id = art.Id
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		dao := NewArticleDao(tx)
		var (
			err error
		)
		if id > 0 {
			err = dao.UpdateById(ctx, art)
		} else {
			id, err = dao.Insert(ctx, art)
		}

		if err != nil {
			return err
		}
		art.Id = id
		pubArt := PublishedArticle(art)
		now := time.Now().UnixMilli()
		pubArt.Ctime = now
		pubArt.Utime = now
		err = tx.Clauses(clause.OnConflict{
			// 对mysql不起效，但可兼容别的方言
			// 下边的代码对应到mysql：INSERT xxx ON DUPLICATE KEY SET ‘title’=？
			// 其它方言：sqlite INSERT xxx ON CONFLICT DO UPDATES WHERE
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":   pubArt.Title,
				"content": pubArt.Content,
				"utime":   now,
			}),
		}).Create(&pubArt).Error
		return err
	})
	return id, err
}

func (a *ArticleGormDao) SyncV1(ctx context.Context, art Article) (int64, error) {
	tx := a.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// 防止后边业务panic
	defer tx.Rollback()
	dao := NewArticleDao(tx)

	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = dao.UpdateById(ctx, art)
	} else {
		id, err = dao.Insert(ctx, art)
	}

	if err != nil {
		return 0, err
	}
	art.Id = id
	pubArt := PublishedArticle(art)
	now := time.Now().UnixMilli()
	pubArt.Ctime = now
	pubArt.Utime = now
	err = tx.Clauses(clause.OnConflict{
		// 对mysql不起效，但可兼容别的方言
		// 下边的代码对应到mysql：INSERT xxx ON DUPLICATE KEY SET ‘title’=？
		// 其它方言：sqlite INSERT xxx ON CONFLICT DO UPDATES WHERE
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   pubArt.Title,
			"content": pubArt.Content,
			"utime":   now,
		}),
	}).Create(&pubArt).Error
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, nil
}

func NewArticleDao(db *gorm.DB) ArticleDao {
	return &ArticleGormDao{
		db: db,
	}
}

func (a *ArticleGormDao) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := a.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (a *ArticleGormDao) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	res := a.db.WithContext(ctx).Model(&Article{}).
		Where("id=? AND author_id=?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("更新失败，ID不对或作者不对")
	}
	return nil
}

type Article struct {
	Id       int64  `gorm:"primaryKey, autoIncrement"`
	Title    string `gorm:"type=varchar(4096)"`
	Content  string `gorm:type=BLOB`
	AuthorId int64  `gorm:"index"`
	Status   int64
	Ctime    int64
	Utime    int64
}

// PublishedArticle 衍生类型
type PublishedArticle Article

// PublishedArticleV1 使用组合的方式
type PublishedArticleV1 struct {
	Article
}
