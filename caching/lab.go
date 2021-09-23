package main

import "fmt"

func main() {
	active := map[int]int{}
	isActive, exists := active[3]
	fmt.Println(isActive, exists)
}
