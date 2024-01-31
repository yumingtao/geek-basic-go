package channel

import (
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	// 声明
	//var ch chan struct{}
	// 声明并创建, 不带buffer
	//ch1 := make(chan int)
	// 声明并创建且带buffer
	ch2 := make(chan int, 3)
	ch2 <- 343
	data := <-ch2
	t.Log(data)
	close(ch2)
}

func TestChannelClose(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 3
	val, ok := <-ch
	t.Log("读到数据了么？", ok, val)
	close(ch)
	// 关闭了channel之后再向channel写入数据会引起panic
	// ch <- 123
	val, ok = <-ch
	t.Log("读到数据了么？", ok, val)
	// 关闭已经关闭的channel会panic
	// close(ch)
}

func TestChannelLoop(t *testing.T) {
	ch := make(chan int, 1)
	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
			time.Sleep(time.Second)
		}
		close(ch)
	}()
	now := time.Now()
	for val := range ch {
		t.Log(val, time.Since(now).Milliseconds())
	}
}

func TestChannelBlocking(t *testing.T) {
	ch := make(chan int)
	b1 := BigStruct{}
	go func() {
		var b BigStruct
		// 设置了一个没有缓冲区的channel，但是没有读取数据的操作，一直blocking
		// 这样b和b1一直不能被回收，会产生goroutine泄露
		ch <- 123
		t.Log(b, b1)
	}()
}

type BigStruct struct {
}

func TestChannelSelect(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 2)
	go func() {
		time.Sleep(time.Second)
		ch1 <- 123
	}()

	go func() {
		time.Sleep(time.Second)
		ch2 <- 234
	}()

	select {
	case val := <-ch1:
		t.Log("ch1来了", val)
	case val := <-ch2:
		t.Log("ch2来了", val)
	}
}
