package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

/*
	两种不同类型的心跳
	1. 在一段时间间隔内发出的心跳
		可以在并发程序中用于报告平安
	2. 在工作单元开始时发出的心跳

*/
func main() {
	test2()
}

// DoWork 一个可以发送心跳的go程 间隔心跳
func DoWork(ctx context.Context, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
	heartBeat := make(chan interface{})
	results := make(chan time.Time)
	go func() {
		defer close(heartBeat)
		defer close(results)
		pulse := time.Tick(pulseInterval) // 心跳
		workGen := time.Tick(2 * pulseInterval)
		sendPulse := func() {
			select {
			case heartBeat <- pulse:
			default: // 可能没有接收者
			}
		}
		sendResult := func(r time.Time) {
			for {
				select {
				case <-ctx.Done(): // 抢占
					return
				case <-pulse:
					sendPulse()
				case results <- r:
					return
				}
			}
		}
		// 真正的控制中心
		for {
			select {
			case <-ctx.Done():
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			}
		}
	}()
	return heartBeat, results
}

func test1() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	const timeout = 2 * time.Second
	heartBeat, results := DoWork(ctx, timeout/2) // 是我们对于超时有额外的响应时间
	for {
		select {
		case _, ok := <-heartBeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Println("result:", r.Second())
		case <-time.After(timeout):
			fmt.Println("timeout")
			return
		}
	}
}

// DoWork2 在工作单元开始时发出的心跳
func DoWork2(ctx context.Context) (<-chan interface{}, <-chan int) {
	heartStream := make(chan interface{})
	workStream := make(chan int)
	go func() {
		defer func() {
			close(heartStream)
			close(workStream)
		}()
		for i := 0; i < 10; i++ {
			select {
			case heartStream <- struct{}{}: // 开始任务时发送信号
			default:
			}
			select {
			case <-ctx.Done():
				return
			case workStream <- rand.Intn(10):
			}
		}
	}()
	return heartStream, workStream
}

func test2() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	heartBeat, results := DoWork2(ctx)
	for {
		select {
		case _, ok := <-heartBeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Println("results:", r)
		}
	}
}
