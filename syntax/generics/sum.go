package main

import "io"

func Sum[T Number](vals []T) T {
	var res T
	for _, v := range vals {
		res = res + v
	}
	return res
}

func Max[T Number](vals []T) T {
	t := vals[0]
	for i := 1; i < len(vals); i++ {
		if t < vals[i] {
			t = vals[i]
		}
	}
	return t
}

func Min[T Number](vals []T) T {
	t := vals[0]
	for i := 1; i < len(vals); i++ {
		if t > vals[i] {
			t = vals[i]
		}
	}
	return t
}

func Find[T any](vals []T, filter func(t T) bool) T {
	for _, v := range vals {
		if filter(v) {
			return v
		}
	}
	var t T
	return t
}

func Insert[T any](idx int, val T, vals []T) []T {
	if idx < 0 || idx > len(vals) {
		panic("idx不合法")
	}
	vals = append(vals, val) //先扩容
	for i := len(vals) - 1; i > idx; i-- {
		if i-1 >= 0 {
			vals[i] = vals[i-1]
		}
	}
	vals[idx] = val
	return vals
}

type Integer int

// 类型约束
type Number interface {
	//~表示衍生类型
	~int | uint | int32
}

func UseSum() {
	res := Sum[int]([]int{1, 2})
	println(res)
	res1 := Sum[Integer]([]Integer{3, 4})
	println(res1)
}

func Closable[T io.Closer]() {
	var t T
	err := t.Close()
	if err != nil {
		return
	}
}
