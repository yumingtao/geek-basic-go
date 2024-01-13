package mongodb

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

type MongoDBTestSuite struct {
	suite.Suite
	col *mongo.Collection
}

func (s *MongoDBTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	t := s.T()
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			fmt.Println(evt.Command)
		},
	}
	opts := options.Client().ApplyURI("mongodb://root:example@localhost:27017").
		SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)
	col := client.Database("webook").Collection("articles")
	s.col = col

	manyRes, err := col.InsertMany(ctx, []any{
		Article{
			Id:       1,
			AuthorId: 123,
		}, Article{
			Id:       2,
			AuthorId: 456,
		},
	})
	assert.NoError(t, err)
	t.Log("插入文档数量：", len(manyRes.InsertedIDs))
}

func (s *MongoDBTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := s.col.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)
	_, err = s.col.Indexes().DropAll(ctx)
	assert.NoError(s.T(), err)
}

func (s *MongoDBTestSuite) TestOr() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	/*filter := bson.M{
		"$or": []bson.M{
			{"authorId": 123},
			{"authorId": 456},
		},
	}
	findRes, err := s.col.Find(ctx, filter)*/
	filter := bson.A{ //切片， value可以是E,D,M
		bson.D{ // E的切片
			bson.E{Key: "id", Value: 1},
		},
		bson.D{ // E的切片
			bson.E{Key: "id", Value: 2},
		},
	}
	findRes, err := s.col.Find(ctx, bson.D{bson.E{Key: "$or", Value: filter}})

	assert.NoError(s.T(), err)
	var arts []Article
	/*for findRes.Next(ctx) {
		var article Article
		err := findRes.Decode(&article)
		assert.NoError(s.T(), err)
		arts = append(arts, article)
	}
	assert.Equal(s.T(), 2, len(arts))
	s.T().Log("查询结果：", arts)*/
	err = findRes.All(ctx, &arts)
	assert.NoError(s.T(), err)
	s.T().Log("查询结果：", arts)
}

func (s *MongoDBTestSuite) TestAnd() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	/*filter := bson.M{
		"$and": []bson.M{
			{"authorId": 123},
			{"authorId": 456},
		},
	}
	findRes, err := s.col.Find(ctx, filter)*/
	filter := bson.A{ //切片， value可以是E,D,M
		bson.D{ // E的切片
			bson.E{Key: "id", Value: 1},
		},
		bson.D{ // E的切片
			bson.E{Key: "authorId", Value: 123},
		},
	}
	findRes, err := s.col.Find(ctx, bson.D{bson.E{Key: "$and", Value: filter}})

	assert.NoError(s.T(), err)
	var arts []Article
	err = findRes.All(ctx, &arts)
	assert.NoError(s.T(), err)
	s.T().Log("查询结果：", arts)
}

func (s *MongoDBTestSuite) TestIn() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	filter := bson.M{
		"authorId": bson.M{
			"$in": []int{
				123,
				456,
			},
		},
	}
	findRes, err := s.col.Find(ctx, filter)
	/*filter := bson.D{bson.E{Key: "id", Value: bson.D{bson.E{Key: "$in", Value: []int{1, 2}}}}}
	findRes, err := s.col.Find(ctx, filter)*/

	assert.NoError(s.T(), err)
	var arts []Article
	err = findRes.All(ctx, &arts)
	assert.NoError(s.T(), err)
	s.T().Log("查询结果：", arts)
}

// 指定查询特定的字段
func (s *MongoDBTestSuite) TestProjection() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	filter := bson.M{
		"authorId": bson.M{
			"$in": []int{
				123,
				456,
			},
		},
	}
	findRes, err := s.col.Find(ctx, filter, options.Find().SetProjection(bson.M{"id": 1}))
	/*filter := bson.D{bson.E{Key: "id", Value: bson.D{bson.E{Key: "$in", Value: []int{1, 2}}}}}
	findRes, err := s.col.Find(ctx, filter,
		options.Find().SetProjection(bson.D{bson.E{Key: "id", Value: 1}}))*/

	assert.NoError(s.T(), err)
	var arts []Article
	err = findRes.All(ctx, &arts)
	assert.NoError(s.T(), err)
	s.T().Log("查询结果：", arts)
}

// 创建索引
func (s *MongoDBTestSuite) TestIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ires, err := s.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{bson.E{"id", 1}},
		Options: options.Index().SetUnique(true).SetName("my_idx_id"),
	})
	assert.NoError(s.T(), err)
	s.T().Log("创建索引：", ires)
}

func TestMongoDBQueries(t *testing.T) {
	suite.Run(t, &MongoDBTestSuite{})
}
