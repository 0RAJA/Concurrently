package main

import (
	"context"
	"fmt"
	"time"
)

/*
	两种不同类型的心跳
	1. 在一段时间间隔内发出的心跳
		可以在并发程序中用于报告平安
	2. 在工作单元开始时发出的心跳

*/
func main() {
	test1()
}

// DoWork 一个可以发送心跳的go程 间隔心跳 模拟滴答声
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
		sendResult := func(r time.Time) { // 发送滴答
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

// DoWork2 在工作单元开始时发出的心跳 防止设置超时导致还没开始工作就超时了
// 根据给定数组内容生成int流
func DoWork2(ctx context.Context, nums ...int) (<-chan interface{}, <-chan int) {
	heartStream := make(chan interface{}, 1) // 至少可以发送一个心跳
	intStream := make(chan int)
	go func() {
		defer func() {
			close(heartStream)
			close(intStream)
		}()
		time.Sleep(time.Second)
		for _, num := range nums {
			select {
			case heartStream <- struct{}{}: // 开始任务时发送信号
			default: // 防止没人接收心跳
			}
			select {
			case <-ctx.Done():
				return
			case intStream <- num:
			}
		}
	}()
	return heartStream, intStream
}

// DoWork3 如果一个迭代会持续很长时间，可以使用间隔心跳保证安全
func DoWork3(ctx context.Context, pulseInterval time.Duration, nums ...int) (<-chan interface{}, <-chan int) {
	heartStream := make(chan interface{}, 1) // 至少可以发送一个心跳
	intStream := make(chan int)
	go func() {
		defer close(heartStream)
		defer close(intStream)

		time.Sleep(2 * time.Second)
		pulse := time.NewTicker(pulseInterval)
		defer pulse.Stop()
	numLoop:
		for _, n := range nums {
			for {
				select {
				case <-ctx.Done():
					return
				case <-pulse.C:
					select {
					case heartStream <- struct{}{}:
					default:
					}
				case intStream <- n:
					continue numLoop
				}
			}
		}
	}()
	return heartStream, intStream
}
