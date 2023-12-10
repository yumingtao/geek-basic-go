package wire

import "fmt"

func UseRepository() {
	// 因为这个文件没有go:build wireinject标签,所以这里调用的是wire_gen里的方法
	repo := InitUserRepository()
	fmt.Println(repo)
}
