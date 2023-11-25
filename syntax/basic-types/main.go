package main

import (
	"math"
	"strconv"
	"unicode/utf8"
)

func main() {
	var a int = 456
	var b int = 123
	println(a + b)
	println(a - b)
	println(a * b)
	println(a / b)
	println(float64(a) / float64(b))
	a++
	println(a)
	b--
	println(b)
	var c float64 = 13.2
	//println(a + c)
	println(c)
	println(math.Abs(-13.2))
	ExtremeNum()
	String()
	Byte()
	Bool()
}

func ExtremeNum() {
	println(math.MinInt64)
	println("float64 最小正数：", math.SmallestNonzeroFloat64)
	println("float32 最小正数：", math.SmallestNonzeroFloat32)
}

func String() {
	//He said "Hello, GO"
	println("Hello, GO!")
	println("He said \"Hello, GO\"")
	println(`Hello, GO!
换行，换行
`)
	println("Hello, " + strconv.Itoa(133))
	println(len("hello"))
	println(len("你好"))
	println(utf8.RuneCountInString("你好"))
}

func Byte() {
	var a byte = 12
	println(a)
	var b byte = 'b'
	println(b)
	var str string = "hello"
	var bs []byte = []byte(str)
	var str1 string = string(bs)
	println(str1)
}

func Bool() {
	var a bool = true
	var b bool = false
	println(a && b)
	println(a || b)
	println(!a)
	// !(a && b) => !a || !b
	println(!(a && b))
	// !(a || b) => !a && !b
	println(!(a || b))
}
