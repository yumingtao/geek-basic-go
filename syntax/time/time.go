package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	fmt.Println(now)
	interval := time.Minute * 3
	result := now.Add(-interval)
	fmt.Println(result)
}
