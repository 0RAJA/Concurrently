package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

/*
	复制请求
	将一个请求复制到多个处理程序，选择最快处理完成的结果
*/

// Work 处理程序
func Work(ctx context.Context, id int, wg *sync.WaitGroup, resultChan chan<- int) {
	defer wg.Done()
	start := time.Now()
	costTime := time.Duration(rand.Intn(5)+1) * time.Second // 模拟耗时
	select {
	case <-ctx.Done():
	case <-time.After(costTime):
		select {
		case <-ctx.Done():
		case resultChan <- id:
		}
	}
	tookTime := time.Since(start)
	if tookTime < costTime {
		tookTime = costTime
	}
	log.Printf("%v cost %v\n", id, tookTime)
}

// DO 开启10个程序
func test1() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	resultChan := make(chan int)
	n := 10
	wg.Add(n)
	for i := 0; i < n; i++ {
		go Work(ctx, i, &wg, resultChan)
	}
	first := <-resultChan
	cancel()
	wg.Wait()
	close(resultChan)
	fmt.Println("received result:", first)
}

func main() {
	test1()
}
