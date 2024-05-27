package cronjob

import (
	"context"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer ticker.Stop()
	// 每秒钟会有一个信号
	for {
		select {
		case <-ctx.Done():
			t.Log("循环结束")
			goto end
			// break 是没有效果的
		case now := <-ticker.C:
			t.Log("过了一秒", now.UnixMilli())
		}
	}
end:
	t.Log("goto 过来了，结束程序")
}
