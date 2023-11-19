package main

func Map() {
	m1 := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	println(m1)
	//和slice一样预估容量
	m2 := make(map[string]string)
	m2["key2"] = "value2"
	println(m2)

	val1, ok := m1["key1"]
	if ok {
		println(val1)
	}
	println(val1, ok)

	val2, ok := m1["key2"]
	//打印"" false
	println(val2, ok)

	val2 = m2["key2"]
	println(val2)

	val2 = m2["key1"]
	//打印""
	println(val2)
	println(len(m2))

	//map遍历是随机的
	for k, v := range m1 {
		println(k, v)
	}

	for k := range m1 {
		println(k, m1[k])
	}

	for _, v := range m1 {
		println(v)
	}

	delete(m1, "key1")
	for k, v := range m1 {
		println(k, v)
	}
	//在switch里，值必须是可比较的
	//在map里，key必须是可比较的
	//在Go里可以比较：Go在编译的时候，运行的时候，能够判断出来元素是不相等的
	//基本类型和string都是可比较的
	//如果元素是可比较的，那么该数组也是可以比较的
	//切片是不可以比较的
}
