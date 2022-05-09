package main

import "fmt"

/*
	约束是一种确保了信息只能从一个并发过程中获取到的简单且强大的方法。
	分为：特定约束(公约)和词法约束
*/

func T1() {
	data := make([]int, 4)
	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}
	handleData := make(chan int)
	go loopData(handleData)
	for num := range handleData {
		fmt.Println(num)
	}
}

func T2() {
	chanOwner := func() <-chan int {
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}
	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Println(result)
		}
		fmt.Println("Done")
	}
	results := chanOwner()
	consumer(results)
}

func main() {
	T2()
}
