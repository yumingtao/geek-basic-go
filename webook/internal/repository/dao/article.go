package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
}

type ArticleGormDao struct {
	db *gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDao {
	return &ArticleGormDao{
		db: db,
	}
}

func (a ArticleGormDao) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := a.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (a ArticleGormDao) UpdateById(ctx context.Context, art Article) error {
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
