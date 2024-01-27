package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InteractiveDao interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	InsertLikeInfo(ctx context.Context, biz string, aid int64, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, aid int64, uid int64) error
	InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
}

type GormInteractiveDao struct {
	db *gorm.DB
}

func NewGormInteractiveDao(db *gorm.DB) InteractiveDao {
	return &GormInteractiveDao{
		db: db,
	}
}

func (dao *GormInteractiveDao) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	now := time.Now().UnixMilli()
	cb.Ctime = now
	cb.Utime = now
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&cb).Error
		if err != nil {
			return err
		}
		return tx.WithContext(ctx).Clauses(clause.OnConflict{DoUpdates: clause.Assignments(map[string]interface{}{
			"collect_cnt": gorm.Expr("`collect_cnt` + 1"),
			"utime":       now,
		}),
		}).Create(&Interactive{
			Biz:        cb.Biz,
			BizId:      cb.BizId,
			CollectCnt: 1,
			Ctime:      now,
			Utime:      now,
		}).Error
	})

	return err
}

func (dao *GormInteractiveDao) InsertLikeInfo(ctx context.Context, biz string, aid int64, uid int64) error {
	now := time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{DoUpdates: clause.Assignments(map[string]interface{}{
			"utime":  now,
			"status": 1,
		}),
		}).Create(&UserLikeBiz{
			Uid:    uid,
			Biz:    biz,
			BizId:  aid,
			Status: 1,
			Ctime:  now,
			Utime:  now,
		}).Error
		if err != nil {
			return err
		}

		return tx.WithContext(ctx).Clauses(clause.OnConflict{DoUpdates: clause.Assignments(map[string]interface{}{
			"like_cnt": gorm.Expr("`like_cnt` + 1"),
			"utime":    now,
		}),
		}).Create(&Interactive{
			Biz:     biz,
			BizId:   aid,
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
		}).Error
	})

	return err
}

func (dao *GormInteractiveDao) DeleteLikeInfo(ctx context.Context, biz string, aid int64, uid int64) error {
	now := time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeBiz{}).
			Where("uid=? AND biz_id=? AND biz=?", uid, aid, biz).
			Updates(map[string]interface{}{
				"utime":  now,
				"status": 0,
			}).Error
		if err != nil {
			return err
		}

		return tx.Model(&Interactive{}).
			Where("biz_id=? AND biz=?", aid, biz).
			Updates(map[string]interface{}{
				"like_cnt": gorm.Expr("`like_cnt` - 1"),
				"utime":    now,
			}).Error
	})

	return err
}

func (dao *GormInteractiveDao) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"read_cnt": gorm.Expr("`read_cnt` + 1"),
			"utime":    now,
		}),
	}).Create(&Interactive{
		BizId:   bizId,
		Biz:     biz,
		ReadCnt: 1,
		Utime:   now,
		Ctime:   now,
	}).Error
}

// Interactive 使用了联合主键<bizId, biz>
type Interactive struct {
	Id         int64  `gorm:"primaryKey,autoincrement"`
	BizId      int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz        string `gorm:"uniqueIndex:biz_type_id;type:varchar(128)"`
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Ctime      int64
	Utime      int64
}

type UserLikeBiz struct {
	Id     int64  `gorm:"primaryKey,autoincrement"`
	Uid    int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	BizId  int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz    string `gorm:"uniqueIndex:uid_biz_type_id;type:varchar(128)"`
	Status int
	Ctime  int64
	Utime  int64
}

type UserCollectionBiz struct {
	Id    int64  `gorm:"primaryKey,autoincrement"`
	Uid   int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	BizId int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz   string `gorm:"uniqueIndex:uid_biz_type_id;type:varchar(128)"`
	Cid   int64  `gorm:"index"` //收藏夹id，本身有索引
	Ctime int64
	Utime int64
}
