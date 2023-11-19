package main

import "geek-basic-go/syntax/variables/access"

var (
	a int     = 12
	b float64 = 33.2
	c int     = 33
	//hh := 7 //只能用于局部且是简单类型
)

const (
	Status0 = iota
	Status1
	Status2
	sss
	Status4 = 5
	Status5
	Status6 = iota
	Status8
)
const (
	MyStatus0 = iota<<4 + 1
	MyStatus1
	MyStatus2
	Mysss
	MyStatus4 = 5
	MyStatus5
)

func main() {
	println(float64(a) + b)
	println(a + c)
	println(Global)
	AccessPackVariables()
	var (
		c int = 3
	)
	println(a + c)
	var d = 1
	println(d)
	var e uint = 222
	println(e)
	var f int
	println(f)
	var g string
	println(g)
	println(access.Global1)
	println(access.Global2)
	//println(access.internal)

	//只能用于局部且是简单类型
	h := 7
	println(h)
	println(access.ConstVariable)
	//access.ConstVariable = "hell0"
	println(Status4)
	println(Status5)
}
