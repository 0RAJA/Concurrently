package main

import (
	"log"
	"os"
	"time"

	"github.com/0RAJA/Concurrently/3_pattern/common"
)

// 治愈异常的goroutine
/*
	当一些goroutine处于异常状态时，尝试对其进行重启
	通过心跳机制判断goroutine的状态（最好在心跳中包含某些信息用于判断其不是火锁）

	管理员:负责监控并重启goroutine
*/

// StartGoroutineFn 创建一个可以监控和重启的goroutine的方式
// 参数：退出channel,管理员的心跳时间
// 返回值：返回管理员心跳的channel
type StartGoroutineFn func(done <-chan any, pulseInterval time.Duration) <-chan any

// NewSteward 新建一个管理员
// 参数：下游的超时时间，创建一个可以监控和重启的goroutine的方式
// 返回值：返回一个创建一个受管理的goroutine和其管理者的函数的创建方式
func NewSteward(timeout time.Duration, startGoroutine StartGoroutineFn) StartGoroutineFn {
	return func(done <-chan any, pulseInterval time.Duration) <-chan any {
		heartBeat := make(chan any)
		go func() {
			defer close(heartBeat)

			var wardDone chan any        // 管理者用于通知下游退出的channel
			var wardHeartbeat <-chan any // 管理员用于接收下游心跳的channel
			startWard := func() {
				wardDone = make(chan any)                                            // 初始化退出channel
				wardHeartbeat = startGoroutine(common.Or(wardDone, done), timeout/2) // 启动下游，其心跳间隔是超时间隔的一半
			}
			startWard()                       // 启动受监管的goroutine
			pulse := time.Tick(pulseInterval) // 定时回复上游的心跳
		monitorLoop:
			for {
				timeoutSignal := time.After(timeout)
				for {
					select {
					case <-pulse: // 回复心跳
						select {
						case heartBeat <- struct{}{}:
						default:
						}
					case <-wardHeartbeat: // 接收到下游的心跳则继续监视
						continue monitorLoop
					case <-timeoutSignal: // 没收到下游的心跳则重启下游
						log.Println("stewart: ward unhealthy;restarting")
						close(wardDone)
						startWard() // 使用之前的方式重新启动下游
						continue monitorLoop
					case <-done:
						return
					}
				}
			}
		}()
		return heartBeat
	}
}

// 不正常的go程
func badWorker() StartGoroutineFn {
	return func(done <-chan any, pulseInterval time.Duration) <-chan any {
		log.Println("ward: Hello, I am irresponsible")
		go func() {
			<-done
			log.Println("ward: I am halting")
		}()
		return nil // 故意阻塞
	}
}

// 受管理的go程:生成int流，可以启动多个管理区副本
func generatorIntStream(done <-chan interface{}, intList ...int) (<-chan interface{}, StartGoroutineFn) {
	intChanStream := make(chan (<-chan interface{}))
	intStream := common.Bridge(done, intChanStream) // 从intChanStream读int流
	return intStream, func(done <-chan any, pulseInterval time.Duration) <-chan any {
		intStream := make(chan interface{}) // 给管理区传递信息
		heartBeat := make(chan interface{})
		go func() {
			defer close(intStream)
			defer log.Println("over")
			log.Println("start")
			select {
			case intChanStream <- intStream: // 尝试塞intStream进intChanStream 塞不进去说明有其他的实例正在工作
			case <-done:
				return
			}
			pulse := time.Tick(pulseInterval)
			for {
			valueLoop:
				// 真正的工作 -> 往intStream里一直塞数据
				for _, intVal := range intList {
					if intVal < 0 {
						log.Printf("negative value:%v\n", intVal)
						return
					}
					time.Sleep(pulseInterval * 2) // 模拟下不正常运行
					for {
						select {
						case <-pulse:
							select {
							case heartBeat <- struct{}{}:
							default:
							}
						case intStream <- intVal:
							continue valueLoop
						case <-done:
							return
						}
					}
				}
			}
		}()
		return heartBeat
	}
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC | log.Lshortfile)

	done := make(chan any)
	time.AfterFunc(time.Minute, func() { // 1min 后退出
		log.Println("main: halting stewart and ward.")
		close(done)
	})
	intStream, startFn := generatorIntStream(done, 1, 2, 3, 4)
	heartBeat := NewSteward(4*time.Second, startFn)(done, 4*time.Second)
	for {
		select {
		case <-done:
			log.Println("Done")
			return
		case <-heartBeat:
			log.Println("Steward is healthy")
		case val := <-intStream:
			log.Println("received intVal:", val.(int))
		}
	}
}
