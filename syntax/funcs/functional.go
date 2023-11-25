package main

func Func6() {
	myFunc3 := Func3
	str, err := myFunc3(1, 4)
	println(str, err)
}

func Func7() {
	//函数内部的匿名方法
	fn := func(name string) string {
		return "Hello," + name
	}
	str := fn("Tony")
	println(str)
}

func Func8() func(name string) string {
	return func(name string) string {
		return "Hello, 方法做为返回值，" + name
	}
}

func Func8Invoke() {
	fn := Func8()
	println(fn("yumingtao"))
}

func Func9() {
	fn := func(name string) string {
		return "Hello, 匿名方法直接发起调用，" + name
	}("yumingtao")
	println(fn)
}
