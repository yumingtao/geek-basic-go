package main

import "fmt"

func main() {
	useDelete()
}

func useDelete() {
	s0 := []int{1, 2, 3, 4, 5, 6}
	s0 = Delete(s0, 0)
	fmt.Println(s0)
	fmt.Printf("s0:%v, len:%d, cap:%d \n", s0, len(s0), cap(s0))

	s0 = Delete(s0, 4)
	fmt.Println(s0)
	fmt.Printf("s0:%v, len:%d, cap:%d \n", s0, len(s0), cap(s0))

	s0 = Delete(s0, 2)
	fmt.Println(s0)
	fmt.Printf("s0:%v, len:%d, cap:%d \n", s0, len(s0), cap(s0))

	s0 = Delete(s0, 1)
	fmt.Println(s0)
	fmt.Printf("s0:%v, len:%d, cap:%d \n", s0, len(s0), cap(s0))

	/*s0 = Delete(s0, 100)
	fmt.Println(s0)
	fmt.Printf("s0:%v, len:%d, cap:%d \n", s0, len(s0), cap(s0))*/
}
