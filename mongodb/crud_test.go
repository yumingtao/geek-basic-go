package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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
	insertRes, err := col.InsertOne(ctx, Article{
		Id:       1,
		Title:    "我的标题",
		Content:  "我的内容",
		AuthorId: 123,
	})
	assert.NoError(t, err)
	oid := insertRes.InsertedID.(primitive.ObjectID)
	// t.Log(string(oid[:12]))
	t.Log("插入ID:", oid)

	//filter := bson.D{bson.E{"id", 1}}
	filter := bson.M{
		"id": 1,
	}
	findRes := col.FindOne(ctx, filter)
	if errors.Is(findRes.Err(), mongo.ErrNoDocuments) {
		t.Log("没找到数据")
	} else {
		assert.NoError(t, findRes.Err())
		var art Article
		err = findRes.Decode(&art)
		assert.NoError(t, err)
		t.Log("找到数据：", art)
	}
}

func TestUpdateOne(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

	// 更新文档
	updateFilter := bson.D{bson.E{Key: "id", Value: 1}}
	// 这种方法没有真正更新到title，而是增加了两列"key"和"value"
	set := bson.D{bson.E{Key: "$set", Value: bson.E{Key: "title", Value: "新的标题"}}}
	// 这种方法可以更新数据
	/*set := bson.D{bson.E{Key: "$set", Value: bson.M{
		"title": "新的标题",
	}}}*/
	// 下边这种方法也可以更新数据
	//updateFilter := bson.M{"id": 1}
	/*set := bson.M{
		"$set": bson.M{
			"title": "新的标题6",
		},
	}*/
	updatedRes, err := col.UpdateOne(ctx, updateFilter, set)
	assert.NoError(t, err)
	t.Log("更新文档数量：", updatedRes.ModifiedCount)
	findRes := col.FindOne(ctx, updateFilter)
	if errors.Is(findRes.Err(), mongo.ErrNoDocuments) {
		t.Log("没找到数据")
	} else {
		assert.NoError(t, findRes.Err())
		var art Article
		err = findRes.Decode(&art)
		assert.NoError(t, err)
		t.Log("找到数据：", art)
	}
}

func TestUpdateMany(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

	// 更新文档
	updateFilter := bson.D{bson.E{Key: "id", Value: 1}}
	updateManyRes, err := col.UpdateMany(ctx, updateFilter, bson.D{bson.E{Key: "$set", Value: Article{
		Content: "新的内容",
	}}})
	assert.NoError(t, err)
	t.Log("更新文档数量：", updateManyRes.ModifiedCount)

	findRes := col.FindOne(ctx, updateFilter)
	if errors.Is(findRes.Err(), mongo.ErrNoDocuments) {
		t.Log("没找到数据")
	} else {
		assert.NoError(t, findRes.Err())
		var art2 Article
		err = findRes.Decode(&art2)
		assert.NoError(t, err)
		t.Log("找到数据：", art2)
	}
}

func TestDelete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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
	deleteFilter := bson.D{bson.E{Key: "id", Value: 1}}
	deleteRes, err := col.DeleteMany(ctx, deleteFilter)
	assert.NoError(t, err)
	t.Log("删除文档数量:", deleteRes.DeletedCount)
	findRes := col.FindOne(ctx, deleteFilter)
	if errors.Is(findRes.Err(), mongo.ErrNoDocuments) {
		t.Log("没找到数据")
	} else {
		assert.NoError(t, findRes.Err())
		var art Article
		err = findRes.Decode(&art)
		assert.NoError(t, err)
		t.Log("找到数据：", art)
	}
}

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}
