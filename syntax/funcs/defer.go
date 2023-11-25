package main

func Defer() {
	defer func() {
		println("defer1, 匿名方法立刻发起了调用!")
	}()
	defer func() {
		println("defer2, 匿名方法立刻发起了调用!")
	}()
}

func DeferClosure() {
	//作为参数传入的，定义defer的时候就确定了
	//作为闭包引入的，执行defer对应闭包的方法时才确定
	i := 0
	//闭包写法
	defer func() {
		println(i)
	}()
	i = 1
}

func DeferClosure1() {
	//作为参数传入的，定义defer的时候就确定了
	//作为闭包引入的，执行defer对应闭包的方法时才确定
	i := 0
	//参数传递
	defer func(i int) {
		println(i)
	}(i)
	/*defer func(j int) {
		println(j)
	}(i)*/
	i = 1
}

func DeferReturn() int {
	a := 0
	defer func() {
		a = 1
	}()
	return a
}

func DeferReturn1() (a int) {
	a = 0
	defer func() {
		a = 1
	}()
	return a
}

func DeferReturn2() *MyStruct {
	a := &MyStruct{
		name: "Jerry",
	}
	//改的是指针指向的值，不是指针a本体
	defer func() {
		a.name = "Tom"
	}()
	return a
}

type MyStruct struct {
	name string
}

func DeferClosureLoop1() {
	//10...10
	for i := 0; i < 10; i++ {
		defer func() {
			println(i)
		}()
	}
}

func DeferClosureLoop2() {
	//9...0
	for i := 0; i < 10; i++ {
		defer func(val int) {
			println(val)
		}(i)
	}
}

func DeferClosureLoop3() {
	//9...0
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			println(j)
		}()
	}
}
