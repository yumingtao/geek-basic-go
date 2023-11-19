package main

import "fmt"

func main() {
	//结构体地址
	u0 := &User{}
	println(u0)

	u1 := &User{
		Age:  33,
		Name: "Tony",
	}
	println(u1)
	fmt.Printf("name:%s, age:%d \n", u1.Name, u1.Age)

	u2 := User{}
	//println(u2)
	//+表示把里边的字段打印出来，v是Go语言内置的形态
	fmt.Printf("u2:%+v \n", u2)
	u2.Name = "Jerry"
	fmt.Printf("u2:%+v \n", u2)

	u3 := new(User)
	fmt.Printf("u3:%+v \n", u3)
	u3.Name = "Tom"
	fmt.Printf("u3:%+v \n", u3)

	var u4 User
	fmt.Printf("u4:%+v \n", u4)
	u4.Name = "yumingtao"
	fmt.Printf("u4:%+v \n", u4)

	var u5 *User
	//u5.Name = "hello"
	//(*u5).Name = "hello"
	fmt.Printf("u5:%+v \n", u5)
	/*var u51 = *u5
	fmt.Printf("u51:%+v \n", u51)*/

	u6 := User{1, "world", "!"}
	fmt.Printf("u6:%+v \n", u6)

	//UseList()
	//ChangeUser()
	//UseFish()
	Components()
}

func UseList() {
	l1 := LinkedList{}
	l1Ptr := &l1
	var l2 = *l1Ptr
	fmt.Printf("l2:%+v \n", l2)
}
