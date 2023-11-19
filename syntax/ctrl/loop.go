package main

import (
	"fmt"
)

func Loop1() {
	for i := 0; i < 10; i++ {
		println(i)
	}

	for i := 0; i < 10; {
		println(i)
		i++
	}
}

func Loop2() {
	i := 0
	for i := 5; i < 10; i++ {
		println(i)
	}
	//相当于while
	for i < 10 {
		println(i)
		i++
	}
}

func LoopForArray() {
	println("遍历数组")
	arr := [3]int{1, 3, 7}
	for i, val := range arr {
		println(i, val)
	}
}

func LoopForSlice() {
	println("遍历切片")
	arr := []int{5, 7}
	for i, val := range arr {
		println(i, val)
	}
}

func LoopForMap() {
	println("遍历Map")
	m := map[string]int{
		"key1": 2,
		"key2": 3,
	}
	for k, v := range m {
		println(k, v)
	}

	for k := range m {
		println(k, m[k])
	}

	for _, v := range m {
		println(v)
	}
}

func LoopBug() {
	users := []User{
		{
			name: "Tom",
		},
		{
			name: "Jerry",
		},
	}
	m := make(map[string]*User, 2)
	for _, v := range users {
		m[v.name] = &v
	}

	for k, v := range m {
		fmt.Printf("name: %s user:%v \n", k, v)
	}

	println("我是一个分割线")
	n := make(map[string]*User, 2)
	for i := 0; i < len(users); i++ {
		n[users[i].name] = &users[i]
	}

	for k, v := range n {
		fmt.Printf("name: %s user:%v \n", k, v)
	}

	println("我是一个分割线")
	k := make(map[string]User, 2)
	for i := 0; i < len(users); i++ {
		k[users[i].name] = users[i]
	}

	for k, v := range k {
		fmt.Printf("name: %s user:%v \n", k, v)
	}

	println("我是一个分割线")
	l := make(map[string]User, 2)
	for _, u := range users {
		//l[u.name] = users[i]
		l[u.name] = u
	}

	for k, v := range l {
		fmt.Printf("name: %s user:%v \n", k, v)
	}
}

type User struct {
	name string
}

func Switch(status int) string {
	switch status {
	case 0:
		return "初始化"
	case 1:
		str := "运行中"
		//return "运行中"
		return str
		//break //不需要break
		/*default://可以没有default
		return "未知状态"*/
	}
	return "未知状态"
}

func Switch2(status int) string {
	// case后边跟bool, 会返回命中的第一个，要保证case后边的条件互斥
	// 建议使用if
	switch {
	case status > 1:
		return "初始化"
	case status > 5:
		str := "运行中"
		//return "运行中"
		return str
		//break //不需要break
		/*default://可以没有default
		return "未知状态"*/
	}
	return "未知状态"
}
