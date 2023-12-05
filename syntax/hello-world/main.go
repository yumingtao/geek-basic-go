package main

import "time"

func main() {
	println(time.Now().UnixMilli())
	now := time.Now().UnixMilli()
	println(now)
	println("Hello, GO!")
	Hello()
}
