package main

import (
	"Concurrently/3_pattern/common"
	"fmt"
)

/*从一系列的channel中消费产生的值 <-chan1 <-chan2 */

func bridge(done <-chan interface{}, chanStream <-chan <-chan any) <-chan any {
	valStream := make(chan any)
	go func() {
		defer close(valStream)
		for {
			var stream <-chan any
			select {
			case mybeStream, ok := <-chanStream: //读取chanStream中的channel
				if !ok {
					return
				}
				stream = mybeStream
			case <-done:
				return
			}
			for val := range common.OrDone(done, stream) { //读取channel内容发送回去
				select {
				case <-done:
					return
				case valStream <- val:
				}
			}
		}
	}()
	return valStream
}

//使用桥接实现一个在一个包含多个channel的channel上实现一个单channel的门面。
//创建10个channel，每个channel写入一个元素，并把这些channel传入桥接函数

func GenVals() <-chan <-chan interface{} {
	chanStream := make(chan (<-chan interface{}))
	go func() {
		defer close(chanStream)
		for i := 0; i < 10; i++ {
			stream := make(chan interface{}, 1)
			stream <- 1
			close(stream)
			chanStream <- stream
		}
	}()
	return chanStream
}

func main() {
	done := make(chan interface{})
	defer close(done)
	for v := range bridge(done, GenVals()) {
		fmt.Print(v, " ")
	}
}
