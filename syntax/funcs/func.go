package main

func main() {
	//Invoke()
	//Func7()
	//Func8Invoke()
	//Func9()
	//ClosureInvoke()
	//YourNameInvoke()
	//Defer()
	//DeferClosure()
	//DeferClosure1()
	//println(DeferReturn())
	//println(DeferReturn1())
	//println(DeferReturn2().name)
	//DeferClosureLoop1()
	//DeferClosureLoop2()
	DeferClosureLoop3()
}

func Invoke() {
	str := Func0("Tony")
	println(str)
	str1, err := Func1(1, 2, 4, "Tony")
	println(str1, err)

	str2, _ := Func2(1, 3)
	println(str2, err)

	_, err3 := Func2(1, 3)
	println(err3)

	//_, _ := Func2(1, 3)
	_, _ = Func2(1, 3)
	Func2(1, 3)
}

func Func0(name string) string {
	return "hello," + name
}

func Func1(a, b, c int, d string) (string, error) {
	return "Hello, World!", nil
}

func Func2(a, b int) (str string, err error) {
	str = "Hello, World！"
	return str, nil
}

func Func3(a, b int) (str string, err error) {
	str = "Hello, World！"
	return "Hello, GO!", nil
}

func Func4(a, b int) (str string, err error) {
	str = "Hello, World！"
	return
}

//返回值要么都有名字，要么都没名字
/*func Fun5(a, b int)  (str string, error){
	return "abc", nil
}*/
