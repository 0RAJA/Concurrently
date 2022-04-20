package main

import (
	"fmt"
	"math/rand"
)

/*
	pipeline 是可以在系统中形成抽象的另一种工具（另一种比如函数，结构体，方法）可以用于流式处理或者批处理数据
	pipeline 是将一系列数据输入执行操作，然后将数据传回的系统。这些操作被称为stage
	pipeline 的属性是：
		一个stage消耗并返回相同的类型
		一个stage必须用语言来表达，以便它可以被传递
*/

//基础的用法
func t1() {
	multiply := func(values []int, multiplier int) []int {
		results := make([]int, len(values))
		for i, v := range values {
			results[i] = v * multiplier
		}
		return results
	}
	add := func(values []int, additive int) []int {
		results := make([]int, len(values))
		for i, v := range values {
			results[i] = v + additive
		}
		return results
	}
	ints := []int{1, 2, 3, 4}
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}
}

type Option interface {
	~int | ~float32 | ~float64
}

func multiply[V Option](done <-chan interface{}, valueStream <-chan V, multiplier V) <-chan V {
	results := make(chan V)
	go func() {
		defer close(results)
		for v := range valueStream {
			select {
			case <-done:
				return
			case results <- v * multiplier:
			}
		}
	}()
	return results
}

func add[V Option](done <-chan interface{}, valueStream <-chan V, additive V) <-chan V {
	results := make(chan V)
	go func() {
		defer close(results)
		for v := range valueStream {
			select {
			case <-done:
				return
			case results <- v + additive:
			}
		}
	}()
	return results
}

func generator[V Option](done <-chan interface{}, values ...V) <-chan V {
	results := make(chan V)
	go func() {
		defer close(results)
		for _, v := range values {
			select {
			case <-done:
				return
			case results <- v:
			}
		}
	}()
	return results
}

//使用channel来构建
func t2() {
	done := make(chan interface{})
	defer close(done)
	valueStraem1 := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, valueStraem1, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}
	valueStraem2 := generator(done, 1.1, 2.2, 3.3, 4.4)
	pipeline2 := multiply(done, add(done, multiply(done, valueStraem2, 2), 1), 2)
	for v := range pipeline2 {
		fmt.Println(v)
	}
}

//简单生成器
func t3() {
	done := make(chan interface{})
	defer close(done)
	for num := range take(done, repeat(done, 1), 10) {
		fmt.Println(num)
	}
}

//重复生成值
func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}

//取出num个数后结束
func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	results := make(chan interface{})
	go func() {
		defer close(results)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case results <- <-valueStream:
			}
		}
	}()
	return results
}

func t4() {
	done := make(chan interface{})
	defer close(done)
	randFn := func() interface{} { return rand.Int() }
	for num := range take(done, repeatFn(done, randFn), 10) {
		fmt.Println(num)
	}
}

//重复调用函数
func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	results := make(chan interface{})
	go func() {
		defer close(results)
		for {
			select {
			case <-done:
				return
			case results <- fn():
			}
		}
	}()
	return results
}

func toString(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {
	stringStream := make(chan string)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(string):
			}
		}
	}()
	return stringStream
}

func main() {
	t4()
}
