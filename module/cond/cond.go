package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

/*
	Cond 让一组goroutine在特定条件下被唤醒
	一个goroutine的集合点，等待或发布一个event
*/

var status int64

//广播唤醒
func broadcast(c *sync.Cond) {
	c.L.Lock()
	defer c.L.Unlock()
	atomic.StoreInt64(&status, 1)
	c.Broadcast() //全部唤醒
}

//监听
func listen(c *sync.Cond) {
	c.L.Lock()
	defer c.L.Unlock()
	for atomic.LoadInt64(&status) != 1 {
		fmt.Println("wait")
		c.Wait() //等待唤醒自动执行 c.L.Unlock()
	}
	fmt.Println("listen")
}

// E1 唤醒多个go程
func E1() {
	c := sync.NewCond(&sync.Mutex{})
	for i := 0; i < 10; i++ {
		go listen(c)
	}
	time.Sleep(100 * time.Microsecond)
	c.Signal() //唤醒一个最早的goroutine
	time.Sleep(100 * time.Microsecond)
	go broadcast(c)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

// E2 模拟队列的增加删除，在至少有一个项目被添加到队列中后再执行下一个项目
func E2() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)
	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Remove from queue")
		c.L.Unlock()
		c.Signal() //唤醒一个等待的go程
	}
	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 { //可能是其他的信号，需要进行检查
			c.Wait()
		}
		fmt.Println("Added queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(time.Second)
		c.L.Unlock()
	}
}

type Button struct {
	Clicked *sync.Cond
}

func E3() {
	button := Button{Clicked: sync.NewCond(new(sync.Mutex))}
	//注册一个函数，允许我们在等待信号后执行对应函数
	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup //为了确保所有注册的程序已经执行在各自的go程中
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			for {
				c.L.Lock()
				c.Wait()
				fn()
				c.L.Unlock()
			}
		}()
		goroutineRunning.Wait()
	}
	var clickRegistered sync.WaitGroup //为了确保所有注册的函数执行完毕
	subscribe(button.Clicked, func() {
		fmt.Println("test1")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("test2")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("test2")
		clickRegistered.Done()
	})
	startButton := func() {
		clickRegistered.Add(3)
		button.Clicked.Broadcast() //注册一个用户按键程序来模拟
		clickRegistered.Wait()     //等待注册的函数执行完毕
	}
	startButton() //唤醒所有注册的函数，然后执行它们
	startButton() //唤醒所有注册的函数，然后执行它们
	startButton() //唤醒所有注册的函数，然后执行它们
}

func main() {
	//E2()
	E3()
}
