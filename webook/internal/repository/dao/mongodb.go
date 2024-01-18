package dao

import (
	"context"
	"errors"
	"geek-basic-go/webook/internal/domain"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var ErrPossibleIncorrectAuthor = errors.New("更新失败，ID不对或作者不对")

type MongoDBArticleDao struct {
	node    *snowflake.Node
	col     *mongo.Collection
	liveCol *mongo.Collection
}

func NewMongoDBArticleDao(mdb *mongo.Database, node *snowflake.Node) *MongoDBArticleDao {
	return &MongoDBArticleDao{
		node:    node,
		col:     mdb.Collection("articles"),
		liveCol: mdb.Collection("published_articles"),
	}
}

func (m *MongoDBArticleDao) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	art.Id = m.node.Generate().Int64()
	_, err := m.col.InsertOne(ctx, &art)
	return art.Id, err
}

func (m *MongoDBArticleDao) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	filter := bson.D{bson.E{Key: "id", Value: art.Id}, bson.E{Key: "author_id", Value: art.AuthorId}}
	set := bson.D{bson.E{Key: "$set", Value: bson.M{
		"title":   art.Title,
		"content": art.Content,
		"status":  art.Status,
		"utime":   now,
	}}}
	res, err := m.col.UpdateOne(ctx, filter, set)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrPossibleIncorrectAuthor
	}
	return nil
}

func (m *MongoDBArticleDao) Sync(ctx context.Context, art Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = m.UpdateById(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	// liveCol insert or update
	now := time.Now().UnixMilli()
	art.Utime = now
	art.Id = id
	filter := bson.D{bson.E{Key: "id", Value: art.Id}, bson.E{Key: "author_id", Value: art.AuthorId}}
	sets := bson.D{bson.E{Key: "$set", Value: art},
		bson.E{Key: "$setOnInsert", Value: bson.D{bson.E{Key: "ctime", Value: now}}}}
	_, err = m.liveCol.UpdateOne(ctx, filter, sets, options.Update().SetUpsert(true))
	return id, err
}

func (m *MongoDBArticleDao) SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error {
	filter := bson.D{bson.E{Key: "id", Value: id}, bson.E{Key: "author_id", Value: uid}}
	sets := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "status", Value: status}}}}
	res, err := m.col.UpdateOne(ctx, filter, sets)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return ErrPossibleIncorrectAuthor
	}
	_, err = m.liveCol.UpdateOne(ctx, filter, sets)
	return err
}
