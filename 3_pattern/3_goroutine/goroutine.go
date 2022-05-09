package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
	防止go程泄露
	终止方式：
		完成工作了
		不可恢复的错误使它不能工作
		被告知终止工作
			通过在父子go程中建立信号，让父go程可以通知子go程
*/

// T1 在channel上接受go 消费者
func T1() {
	dowork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		results := make(chan interface{})
		go func() {
			defer func() {
				fmt.Println("dowork done")
				close(results)
			}()
			for {
				select {
				case <-done:
					return
				case s, ok := <-strings:
					if ok {
						fmt.Println(s)
					}
				}
			}
		}()
		return results
	}
	done := make(chan interface{})
	terminated := dowork(done, nil)
	go func() {
		<-time.NewTimer(time.Second).C
		fmt.Println("Cancel dowork")
		close(done)
	}()
	<-terminated
	fmt.Println("Done.")
}

// T2 生产者
func T2() {
	newRandStream := func(done <-chan struct{}) <-chan int {
		results := make(chan int)
		go func() {
			defer func() {
				fmt.Println("newRandStream exited")
				close(results)
			}()
			for {
				select {
				case results <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return results
	}
	done := make(chan struct{})
	randStream := newRandStream(done)
	for i := 0; i < 3; i++ {
		fmt.Println(i, ":", <-randStream)
	}
	close(done)
}

func main() {
	T2()
}
