package cronjob

import (
	cron "github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCronExpr(t *testing.T) {
	expr := cron.New(cron.WithSeconds())
	// cron 表达式可以参考https://help.aliyun.com/document_detail/133509.html
	id, err := expr.AddFunc("@every 1s", func() {
		t.Log("我执行了")
	})
	assert.NoError(t, err)
	t.Log("任务", id)
	expr.Start()
	time.Sleep(time.Second * 10)
	// 不要调度新任务，正在执行的继续执行到结束
	ctx := expr.Stop()
	t.Log("发出了停止信号")
	// 整个任务都执行完
	<-ctx.Done()
	// 彻底停下来
	t.Log("彻底停下来，没有任务在执行")
}

type JobFunc func()

func (j JobFunc) Run() {
	j()
}
