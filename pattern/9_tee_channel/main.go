package main

import (
	"Concurrently/pattern/common"
	"fmt"
)

/*
	分割一个来自channel的值，以便将它们发送到代码的两个独立的地方。：
		如一个传递用户指令的channel，我可以将它们发送给指定的执行器以及相应的日志记录下来。
*/

func tee(done <-chan interface{}, in <-chan any) (_, _ <-chan any) {
	out1 := make(chan any)
	out2 := make(chan any)
	go func() {
		defer close(out1)
		defer close(out2)
		for v := range common.OrDone(done, in) {
			var out1, out2 = out1, out2 //本地版本，隐藏外界变量
			for i := 0; i < 2; i++ {    //为了确保两个channel都可以被写入我们使用两次写入
				select {
				case <-done:
					return
				case out1 <- v:
					out1 = nil //同时写入后关闭副本channel来阻塞防止二次写入
				case out2 <- v:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}

func t1() {
	done := make(chan interface{})
	defer close(done)
	//for v := range common.Take(done, common.Repeat(done, 1, 2), 4) {
	//	fmt.Println(v)
	out1, out2 := tee(done, common.Take(done, common.Repeat(done, 1, 2), 4))
	for val1 := range out1 {
		fmt.Println("out1:", val1, "out2:", <-out2)
	}
}

func main() {
	t1()
}
