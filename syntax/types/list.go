package main

/*
*
接口是一组行为的抽象
业务开发也是面向接口编程
*/
type List interface {
	//接口方法不需要些func
	Add(index int, val any)
	Append(val any) error
	Delete(index int) error
}

// 在接口上，快捷键ctrl+i，实现所有方法
func (l *LinkedList) Append(val any) error {
	//TODO implement me
	panic("implement me")
}

func (l *LinkedList) Delete(index int) error {
	//TODO implement me
	panic("implement me")
}

type LinkedList struct {
	head node
	//Head node
}

func (l *LinkedList) Add(index int, val any) {
	//实现这个方法
}

type node struct {
	//next node
	//结构体自引用，要用指针
	next *node
}

func UseList2() {
	l1 := &LinkedList{}
	l1.Add(1, 123)
	l1.Add(2, "dd")
	l1.Add(3, nil)
}
