package context

import (
	"context"
	"testing"
	"time"
)

type key struct {
	key string
}

func TestContextValue1(t *testing.T) {
	ctx := context.WithValue(context.Background(), key{key: "key_struct"}, "value_struct")
	val, ok := ctx.Value(key{key: "key_struct"}).(string)
	t.Log(val, ok)
}
func TestContextValue(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key1", "value1")
	val, ok := ctx.Value("key1").(string)
	t.Log(val, ok)
}

func TestContextCancel(t *testing.T) {
	//context.TODO()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second)
		t.Log("调用cancel")
		cancel()
	}()
	<-ctx.Done()
	t.Log("已经cancel")
	t.Log(ctx.Err())
	/*select {
	case <-ctx.Done():
	case xxx:
		// 业务逻辑
	}*/
}

func TestContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	<-ctx.Done()
	t.Log("超时了")
	t.Log(ctx.Err())
}

func TestContextPatentCancel(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second, func() {
		cancel()
	})
	son, sonCancel := context.WithCancel(parent)
	<-son.Done()
	t.Log("son已经过来了")
	sonCancel()
}

func TestContextPatentCancel1(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	_, sonCancel := context.WithCancel(parent)
	time.AfterFunc(time.Second, func() {
		sonCancel()
	})
	<-parent.Done()
	t.Log("parent已经过来了")
	cancel()
}

func TestContextPatentValue(t *testing.T) {
	parent := context.WithValue(context.Background(), "key1", "value1")
	son2 := context.WithValue(parent, "key1", "son1_value1")
	t.Log(son2.Value("key1"))
	_ = context.WithValue(parent, "key2", "son2_value2")
	t.Log(parent.Value("key2"))
}
