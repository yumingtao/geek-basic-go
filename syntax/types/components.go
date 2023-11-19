package main

import "io"

type NameI interface {
	Name() string
}

func (i Inner) Name() string {
	println("I am Inner, I implement NameI!")
	return "inner"
}

func (i Outer) Name() string {
	println("I am Outer, I implement NameI!")
	return "outer"
}

// 组合不是继承，没有多态
type Outer struct {
	Inner  //组合，如果Inner上有什么方法，Outer是可以直接用的
	Outer2 //组合多个
}

type Outer1 struct {
	*Inner //组合，如果Inner上有什么方法，Outer是可以直接用的
}

type Outer2 struct {
	io.Closer //组合接口
}

type Inner struct {
}

func (i Inner) Hello() {
	//组合使用的是自己的
	println("Hello, 我是Inner。" + i.Name())
}

func Components() {
	var o Outer
	println("调用Name方法，outer")
	o.Hello()
	//o.Name()

	//var o1 Outer1
	//Outer1里边组合的是指针，要初始化
	o1 := Outer1{
		Inner: &Inner{},
	}
	//可以把Inner当成普通的字段初始化
	o1.Inner = &Inner{}
	println("调用Name方法，outer1")
	o1.Hello()
	//o1.Name()
}
