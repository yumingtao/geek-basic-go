package main

import "fmt"

func ChangeUser() {
	println("我是Change User")
	u1 := User{Name: "Tom", Age: 13}
	fmt.Printf("u1:%+v \n", u1)
	fmt.Printf("u1 address %p \n", &u1)

	//当u1.ChangeName的时候，u1放生了复制，值传递
	//相当于ChangeName(u User, name string)
	u1.ChangeName("Jerry")
	u1.ChangeAge(18)
	fmt.Printf("u1:%+v \n", u1)

	println("============u2===============")

	u2 := &User{Name: "yumingtao", Age: 13}
	fmt.Printf("u2: %+v \n", u2)
	fmt.Printf("u2 address %p \n", &u2)

	//复制了一份指针
	u2.ChangeName("Jerry")
	u2.ChangeAge(35)
	fmt.Printf("u2:%+v \n", u2)
}

type User struct {
	Age      int
	Name     string
	NickName string
}

// 使用结构体接收器
func (u User) ChangeName(name string) {
	fmt.Printf("ChangeName u address %p \n", &u)
	u.NickName = name
}

func ChangeName(u User, name string) {
	fmt.Printf("ChangeName u address %p \n", &u)
	u.NickName = name
}

// 使用指针接收器
func (u *User) ChangeAge(age int) {
	fmt.Printf("ChangeAge u address %p \n", &u)
	u.Age = age
}
