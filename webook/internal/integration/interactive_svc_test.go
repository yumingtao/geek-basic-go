package integration

import (
	"context"
	"geek-basic-go/webook/internal/integration/startup"
	"geek-basic-go/webook/internal/repository/dao"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type InteractiveTestSuite struct {
	suite.Suite
	db  *gorm.DB
	rdb redis.Cmdable
}

func (s *InteractiveTestSuite) SetupSuite() {
	s.db = startup.InitDB()
	s.rdb = startup.InitRedis()
}

func (s *InteractiveTestSuite) TearDownSuite() {
	err := s.db.Exec("TRUNCATE TABLE `interactives`").Error
	assert.NoError(s.T(), err)
	err = s.db.Exec("TRUNCATE TABLE `user_like_bizs`").Error
	assert.NoError(s.T(), err)
}

func TestInteractiveService(t *testing.T) {
	suite.Run(t, new(InteractiveTestSuite))
}

func (s *InteractiveTestSuite) TestIncrReadCnt() {
	testCases := []struct {
		name    string
		before  func(t *testing.T)
		after   func(t *testing.T)
		biz     string
		bizId   int64
		wantErr error
	}{
		{
			name: "增加成功，db和redis",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.WithContext(ctx).Create(dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					Ctime:      6,
					Utime:      7,
				}).Error
				assert.NoError(t, err)
				err = s.rdb.HSet(ctx, "interactive:test:2", "read_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Interactive
				err := s.db.WithContext(ctx).Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 7)
				data.Utime = 0
				assert.Equal(t, dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    4, //加了1
					CollectCnt: 4,
					LikeCnt:    5,
					Ctime:      6,
				}, data)
				cnt, err := s.rdb.HGet(ctx, "interactive:test:2", "read_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 4, cnt)
				err = s.rdb.Del(ctx, "interactive:test:2").Err()
				assert.NoError(t, err)
			},
			biz:     "test",
			bizId:   2,
			wantErr: nil,
		},
		{
			// DB成功，缓存没成功
			name: "增加成功，db成功",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.WithContext(ctx).Create(dao.Interactive{
					Id:         2,
					Biz:        "test",
					BizId:      3,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					Ctime:      6,
					Utime:      7,
				}).Error
				assert.NoError(t, err)
				//err = s.rdb.HSet(ctx, "interactive:test:3", "read_cnt", 3).Err()
				//assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Interactive
				err := s.db.WithContext(ctx).Where("id = ?", 2).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 7)
				data.Utime = 0
				assert.Equal(t, dao.Interactive{
					Id:         2,
					Biz:        "test",
					BizId:      3,
					ReadCnt:    4, //加了1
					CollectCnt: 4,
					LikeCnt:    5,
					Ctime:      6,
				}, data)
				cnt, err := s.rdb.Exists(ctx, "interactive:test:3").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
				err = s.rdb.Del(ctx, "interactive:test:2").Err()
				assert.NoError(t, err)
			},
			biz:     "test",
			bizId:   3,
			wantErr: nil,
		},
		{
			name:   "增加成功,都没有",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Interactive
				err := s.db.WithContext(ctx).Where("biz_id = ? AND biz = ?", 4, "test").First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Ctime > 0)
				assert.True(t, data.Utime > 0)
				assert.True(t, data.Id > 0)
				data.Ctime = 0
				data.Utime = 0
				data.Id = 0
				assert.Equal(t, dao.Interactive{
					Biz:     "test",
					BizId:   4,
					ReadCnt: 1,
				}, data)
				// 这里的意思：因为lua脚本判断如果存在key就加1，如果不存在什么都不做
				// 而我们在before函数中没有做任何事，所以redis中是没有这个ke的，所以拿到的cnt是0
				cnt, err := s.rdb.Exists(ctx, "interactive:test:4").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
				err = s.rdb.Del(ctx, "interactive:test:4").Err()
				assert.NoError(t, err)
			},
			biz:     "test",
			bizId:   4,
			wantErr: nil,
		},
	}

	svc := startup.InitInteractiveService()
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			err := svc.IncrReadCnt(context.Background(), tc.biz, tc.bizId)
			assert.Equal(t, tc.wantErr, err)
			tc.after(t)
		})
	}
}

func (s *InteractiveTestSuite) TestLike() {
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64
		uid   int64

		wantErr error
	}{
		{
			name: "点赞，db和redis都有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.WithContext(ctx).Create(dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					Ctime:      6,
					Utime:      7,
				}).Error
				assert.NoError(t, err)
				err = s.rdb.HSet(ctx, "interactive:test:2", "like_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Interactive
				err := s.db.WithContext(ctx).Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 7)
				data.Utime = 0
				assert.Equal(t, dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    6,
					Ctime:      6,
				}, data)
				var likeBiz dao.UserLikeBiz
				err = s.db.WithContext(ctx).
					Where("uid=? AND biz_id=? AND biz=?", 123, 2, "test").
					First(&likeBiz).Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.Id > 0)
				assert.True(t, likeBiz.Ctime > 0)
				assert.True(t, likeBiz.Utime > 0)
				likeBiz.Utime = 0
				likeBiz.Ctime = 0
				likeBiz.Id = 0
				assert.Equal(t, dao.UserLikeBiz{
					Uid:    123,
					Biz:    "test",
					BizId:  2,
					Status: 1,
				}, likeBiz)
				cnt, err := s.rdb.HGet(ctx, "interactive:test:2", "like_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 4, cnt)
				err = s.rdb.Del(ctx, "interactive:test:2").Err()
				assert.NoError(t, err)
			},
			biz:     "test",
			bizId:   2,
			uid:     123,
			wantErr: nil,
		},
		{
			name:   "点赞，db和redis都没有",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Interactive
				err := s.db.WithContext(ctx).Where("biz=? AND biz_id=?", "test", 3).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Id > 0)
				assert.True(t, data.Ctime > 0)
				assert.True(t, data.Utime > 0)
				data.Id = 0
				data.Ctime = 0
				data.Utime = 0
				assert.Equal(t, dao.Interactive{
					Biz:        "test",
					BizId:      3,
					ReadCnt:    0,
					CollectCnt: 0,
					LikeCnt:    1,
				}, data)
				var likeBiz dao.UserLikeBiz
				err = s.db.WithContext(ctx).
					Where("uid=? AND biz_id=? AND biz=?", 124, 3, "test").
					First(&likeBiz).Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.Id > 0)
				assert.True(t, likeBiz.Ctime > 0)
				assert.True(t, likeBiz.Utime > 0)
				likeBiz.Utime = 0
				likeBiz.Ctime = 0
				likeBiz.Id = 0
				assert.Equal(t, dao.UserLikeBiz{
					Uid:    124,
					Biz:    "test",
					BizId:  3,
					Status: 1,
				}, likeBiz)
				cnt, err := s.rdb.Exists(ctx, "interactive:test:3").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
			},
			biz:     "test",
			bizId:   3,
			uid:     124,
			wantErr: nil,
		},
	}

	svc := startup.InitInteractiveService()
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			err := svc.Like(context.Background(), tc.biz, tc.bizId, tc.uid)
			assert.Equal(t, tc.wantErr, err)
			tc.after(t)
		})
	}
}

func (s *InteractiveTestSuite) TestDislike() {
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64
		uid   int64

		wantErr error
	}{
		{
			name: "取消点赞，db和redis都有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.WithContext(ctx).Create(dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					Ctime:      6,
					Utime:      7,
				}).Error
				assert.NoError(t, err)
				err = s.db.WithContext(ctx).Create(dao.UserLikeBiz{
					Id:     1,
					Biz:    "test",
					BizId:  2,
					Uid:    123,
					Status: 1,
					Ctime:  6,
					Utime:  7,
				}).Error
				assert.NoError(t, err)
				err = s.rdb.HSet(ctx, "interactive:test:2", "like_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Interactive
				err := s.db.WithContext(ctx).Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 7)
				data.Utime = 0
				assert.Equal(t, dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    4,
					Ctime:      6,
				}, data)
				var likeBiz dao.UserLikeBiz
				err = s.db.WithContext(ctx).
					Where("uid=? AND biz_id=? AND biz=?", 123, 2, "test").
					First(&likeBiz).Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.Id > 0)
				assert.True(t, likeBiz.Ctime > 0)
				assert.True(t, likeBiz.Utime > 0)
				likeBiz.Utime = 0
				likeBiz.Ctime = 0
				likeBiz.Id = 0
				assert.Equal(t, dao.UserLikeBiz{
					Uid:    123,
					Biz:    "test",
					BizId:  2,
					Status: 0,
				}, likeBiz)
				cnt, err := s.rdb.HGet(ctx, "interactive:test:2", "like_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 2, cnt)
				err = s.rdb.Del(ctx, "interactive:test:2").Err()
				assert.NoError(t, err)
			},
			biz:     "test",
			bizId:   2,
			uid:     123,
			wantErr: nil,
		},
	}

	svc := startup.InitInteractiveService()
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			err := svc.CancelLike(context.Background(), tc.biz, tc.bizId, tc.uid)
			assert.Equal(t, tc.wantErr, err)
			tc.after(t)
		})
	}
}
