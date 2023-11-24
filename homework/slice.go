package main

//Delete
/*
1.不考虑顺序，采用将最后一个元素放到删除位置，然后返回子切片
2.当切片长度小于容量的一半时考虑缩容，重新申请内存
*/
func Delete[T any](vals []T, idx int) []T {
	l := len(vals)
	if idx < 0 || idx >= l {
		panic("idx不合法")
	}
	//保持原来的顺序，从idx开始，后边的元素依次往前挪一位
	/*for i := idx; i < l-1; i++ {
		vals[i] = vals[i+1]
	}*/
	//如果不考虑顺序，直接将最后一个元素放到idx
	vals[idx] = vals[l-1]
	l--
	if 2*l < cap(vals) {
		println("缩容了...")
		newVals := make([]T, l)
		for i := 0; i < l; i++ {
			newVals[i] = vals[i]
		}
		clear(vals)
		return newVals
	}
	return vals[:l]
}

//Delete1
/*
*
没有返回值的方式，每次都重新分配内存
*/
func Delete1[T any](idx int, vals *[]T) {
	l := len(*vals)
	if idx < 0 || idx >= l {
		panic("idx不合法")
	}
	l--
	newVals := make([]T, l)
	d := false
	for i := 0; i < l; i++ {
		if idx == i {
			newVals[i] = (*vals)[i+1]
			d = true
		} else {
			if d {
				newVals[i] = (*vals)[i+1]
			} else {
				newVals[i] = (*vals)[i]
			}
		}
	}
	*vals = newVals
}
