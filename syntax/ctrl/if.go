package main

func IfOnly(age int) string {
	if age >= 18 {
		return "已经成年了"
	}
	return "还是一个孩子"
}

func IfElse(age int) string {
	if age >= 18 {
		return "已经成年了"
	} else {
		return "还是你一个孩子"
	}
}

func IfElseIf(age int) string {
	if age >= 18 {
		return "已经成年了"
	} else if age > 12 {
		return "还是一个少年"
	} else {
		return "还是你一个孩子"
	}
}
func IfElseIf1(age int) string {
	if age >= 18 {
		return "已经成年了"
	} else if age > 12 {
		return "还是一个少年"
	}
	return "还是你一个孩子"
}

func IfNewVariable(start int, end int) string {
	if distance := end - start; distance > 100 {
		println(distance)
		return "太远了"
	} else {
		println(distance)
		return "OK"
	}
	//println(distance)
}
