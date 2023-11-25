package main

func Closure(name string) func() string {
	//闭包
	return func() string {
		return "Hello, 我是一个闭包，" + name
	}
}

func ClosureInvoke() {
	c := Closure("Tony")
	println(c())
}
