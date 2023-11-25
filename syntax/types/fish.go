package main

//Fish
/*
结构体实现接口：鸭子类型
"当看到的东西走起来像鸭子，游泳起来像鸭子，叫起来像鸭子，那么这个东西就可以被称为鸭子"
当一个结构体具备的接口的所有的方法的时候，它就实现了这个接口
*/
type Fish struct {
}

func (f Fish) Swim() {
	println("会游泳的真鱼")
}

// type TypeA TypeB
// 衍生类型，是一个全新的类型
// TypeB 实现了某个接口，不代表TypeA也实现了这个接口
// TypeA 和 TypeB可以互相转化
type FakeFish Fish

type Yu = Fish //注意这个是别名

func (f FakeFish) FakeSwim() {
	println("会游泳的假鱼")
}

func UseFish() {
	f1 := Fish{}
	f1.Swim()

	f2 := FakeFish{}
	//f2.FakeSwim()
	f2.FakeSwim()

	//类型转换
	f3 := Fish(f2)
	f3.Swim()

	println("这里是一个别名")
	yu := Yu{}
	yu.Swim()
}
