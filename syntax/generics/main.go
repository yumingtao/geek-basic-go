package main

import "fmt"

func main() {
	UseSum()
	println(Max([]int{1, 4, 8}))
	println(Min([]int{1, 4, 8}))
	res0 := Insert[int](0, 12, []int{1, 2})
	fmt.Printf("res0:%v \n", res0)

	res1 := Insert[int](2, 12, []int{1, 2})
	fmt.Printf("res1:%v \n", res1)

	res2 := Insert[int](1, 12, []int{1, 2})
	fmt.Printf("res2:%v \n", res2)

	res3 := Insert[int](3, 12, []int{1, 2})
	fmt.Printf("res3:%v \n", res3)
}
