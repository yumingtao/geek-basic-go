package main

type ListV1[T any] interface {
	Add(index int, val T)
	Append(val any) error
	Delete(index int) error
}

func (l LinkedListV1[T]) Add(index int, val T) {
}

func (l LinkedListV1[T]) Append(val any) error {
	//TODO implement me
	panic("implement me")
}

func (l LinkedListV1[T]) Delete(index int) error {
	//TODO implement me
	panic("implement me")
}

type LinkedListV1[T any] struct {
	head *nodeV1[T]
}

type nodeV1[T any] struct {
	data T
}

func UseListV1() {
	l := &LinkedListV1[int]{}
	l.Add(1, 123)
	//l.Add(3, "13") //会有泛型检查，编译报错
}
