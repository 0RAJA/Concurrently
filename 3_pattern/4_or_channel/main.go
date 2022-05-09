package main

import (
	"fmt"
	"time"
)

/*
	希望将一个或者多个完成的channel合并到一个完成的channel，此channel在任何组件channel关闭时关闭
	则可以使用or_channel将这些channel组合起来
*/
//监听多个channel，只要有一个channel有通知就退出,可以将任意数量的channel组合到单个channel中
func T1() {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}
		orDone := make(chan interface{})
		go func() {
			defer close(orDone)
			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...): //递归退出
				}
			}
		}()
		return orDone
	}
	//定时退出通知
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	//等待多个channel中有一个ok
	<-or(
		sig(time.Second),
		sig(time.Second*2),
	)
	fmt.Println(time.Since(start))
}

func main() {
	T1()
}
