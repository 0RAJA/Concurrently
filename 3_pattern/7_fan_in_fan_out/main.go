package main

import (
	"Concurrently/3_pattern/common"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

/*
	有时候pipeline中的各个stage可能在计算上十分昂贵，上游stage可能会被阻塞，我们可以重复使用pipeline的各个stage。在多个go程中重用pipeline的单个stage来并行化来自上游stage的pull
	扇出：启动多个go程来处理来自pipeline的输入过程
	扇入：将多个结果组合到一个channel
		适用：此stage不依赖之前stage计算的结果，运行需要很长时间
*/

func main() {
	t2()
}

func t2() {
	randFn := func() interface{} { return rand.Int63() }
	done := make(chan interface{})
	defer close(done)
	start := time.Now()
	randIntStream := common.ToType[int64](done, common.RepeatFn(done, randFn)) //将生成的数据显式类型转换
	fmt.Println("Prime")
	finders := make([]<-chan interface{}, runtime.NumCPU())
	for i := 0; i < len(finders); i++ {
		finders[i] = common.PrimeFinder(done, randIntStream)
	}
	for prime := range common.ToType[int64](done, common.Take(done, common.FanIn(done, finders), 10)) {
		fmt.Printf("%#v\n", prime)
	}
	fmt.Println(time.Since(start)) //
}

func t1() {
	randFn := func() interface{} { return rand.Int() }
	done := make(chan interface{})
	defer close(done)
	start := time.Now()
	randIntStream := common.ToType[int](done, common.RepeatFn(done, randFn))
	fmt.Println("Prime")
	for prime := range common.Take(done, common.PrimeFinder(done, randIntStream), 10) {
		fmt.Println(prime)
	}
	fmt.Println(time.Since(start)) //58.394557933s 需要很长时间才能找到这些数
}
