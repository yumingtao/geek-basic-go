package dao

import (
	"context"
	"errors"
	"geek-basic-go/webook/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error)
	GetById(ctx context.Context, id int64) (Article, error)
	GetPubById(ctx context.Context, id int64) (PublishedArticle, error)
	GetPubByIds(ctx context.Context, ids []int64) ([]Article, error)
}

type ArticleGormDao struct {
	db *gorm.DB
}

func (a *ArticleGormDao) GetPubByIds(ctx context.Context, ids []int64) ([]Article, error) {
	var res []Article
	err := a.db.WithContext(ctx).Where("id IN ?", ids).Find(&res).Error
	return res, err
}

func (a *ArticleGormDao) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	var res PublishedArticle
	err := a.db.WithContext(ctx).Where("id=?", id).First(&res).Error
	return res, err
}

func (a *ArticleGormDao) GetById(ctx context.Context, id int64) (Article, error) {
	var art Article
	err := a.db.WithContext(ctx).Where("id=?", id).First(&art).Error
	return art, err
}

func (a *ArticleGormDao) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error) {
	var arts []Article
	err := a.db.WithContext(ctx).Where("author_id=?", uid).
		Offset(offset).
		Limit(limit).
		Order("utime DESC").
		Find(&arts).Error
	return arts, err
}

func (a *ArticleGormDao) SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error {
	now := time.Now().UnixMilli()
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id=? and author_id=?", id, uid).
			Updates(map[string]any{
				"utime":  now,
				"status": status,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("更新失败，ID不对或作者不对")
		}

		return tx.Model(&PublishedArticle{}).
			Where("id=?", id).
			Updates(map[string]any{
				"utime":  now,
				"status": status,
			}).Error
	})
	return err
}

func (a *ArticleGormDao) Sync(ctx context.Context, art Article) (int64, error) {
	var id = art.Id
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		dao := NewGormDBArticleDao(tx)
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
	dao := NewGormDBArticleDao(tx)

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
			"status":  pubArt.Status,
			"utime":   now,
		}),
	}).Create(&pubArt).Error
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, nil
}

func NewGormDBArticleDao(db *gorm.DB) ArticleDao {
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
			"status":  art.Status,
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
	Id       int64  `gorm:"primaryKey, autoIncrement" bson:"id,omitempty"`
	Title    string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	Content  string `gorm:"type=BLOB" bson:"content,omitempty"`
	AuthorId int64  `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}

// PublishedArticle 衍生类型
type PublishedArticle Article

// PublishedArticleV1 使用组合的方式
type PublishedArticleV1 struct {
	Article
}
