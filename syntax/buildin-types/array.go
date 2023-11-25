package main

import "fmt"

func Array() {
	println("我是数组啊")
	//三元素数组，直接初始化，大括号元素只能少不能多
	//%v表示打印一个value或是结构体
	a1 := [3]int{1, 3, 4}
	fmt.Printf("a1: %v, len: %d, cap: %d \n", a1, len(a1), cap(a1))

	//少了，给的是类型的默认值
	a2 := [3]int{1, 3}
	fmt.Printf("a2: %v, len: %d, cap: %d \n", a2, len(a2), cap(a2))
	//以下报错
	//var a3 = [3]int
	//a3 := [3]int
	var a3 [3]int
	fmt.Printf("a3: %v, len: %d, cap: %d \n", a3, len(a3), cap(a3))

	//数组不支持append
	//a3 = append(a3, 1)

	fmt.Printf("a[1]: %d \n", a1[1])
	//arr1(100)
}

func arr1(idx int) {
	a1 := [3]int{1, 3, 4}
	fmt.Printf("a[1]: %d", a1[idx])
}

/*
*
弱化版ArrayList，没有add，delete方法
支持子切片，共享底层数组
*/
func Slice() {
	println("我是切片")
	s1 := []int{1, 2, 3, 45}
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))

	s2 := make([]int, 3, 4) //初始化3个元素，容量是4
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))
	s2 = append(s2, 4) //追加一个元素，没有扩容
	fmt.Printf("追加没有扩容 s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))
	s2 = append(s2, 5) //追加一个元素，扩容了
	fmt.Printf("追加扩容了s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	s3 := make([]int, 4) //创建一个容量是4的切片
	fmt.Printf("s3: %v, len: %d, cap: %d \n", s3, len(s3), cap(s3))

	fmt.Printf("s3[3]: %d \n", s3[3])
	//fmt.Printf("s3[2]: %d \n", s3[30])
}

func SubSlice() {
	println("我是子切片")
	s1 := []int{2, 4, 6, 8, 10}
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
	s2 := s1[1:3] //不包含3
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	s3 := s1[2:] //创建[2, len(s1))
	fmt.Printf("s3: %v, len: %d, cap: %d \n", s3, len(s3), cap(s3))

	s4 := s1[:3] //获取[0,end)
	fmt.Printf("s4: %v, len: %d, cap: %d \n", s4, len(s4), cap(s4))
}

/*
*
切片或子切片中的一个发生扩容，就不共享底层数组了
*/
func SharedSlice() {
	println("我是共享切片")
	s1 := []int{1, 2, 3, 4}
	s2 := s1[2:]
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	s2[0] = 99 //s2[0]是s1[2]
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	s2 = append(s2, 199)
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	s2[1] = 1999 //如果没扩容s2[1]是s1[3]
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))
}
