package main

import (
	"fmt"
	"sync"
)

/*
	once 确保在不同的go程中只会调用一次Do()
*/

func main() {
	test1 := func() {
		fmt.Println("test1")
	}
	test2 := func() {
		fmt.Println("test2")
	}
	once := new(sync.Once)
	once.Do(test1)
	once.Do(test2) //不会执行
}
