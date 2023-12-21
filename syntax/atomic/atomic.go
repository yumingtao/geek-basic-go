package main

import "sync/atomic"

func main() {
	var val int32 = 12
	// 原子读，不会读到修改了一半的数据
	val = atomic.LoadInt32(&val)
	println(val)
	// 原子写，即便别的Goroutine在别的cpu核上修改，也能立刻看到
	atomic.StoreInt32(&val, 13)
	// 原子加，返回新值
	newVal := atomic.AddInt32(&val, 1)
	println(newVal)
	// CAS
	// 如果原来的值是13，就修改为14
	swapped := atomic.CompareAndSwapInt32(&val, 13, 15)
	println(swapped)
	swapped = atomic.CompareAndSwapInt32(&val, 14, 15)
	println(swapped)
}
