package common

import "sync"

// Bridge 通过接受传输chan的chan，将值传递给给回去(这个是按顺序读完一个channel才会选择下一个channel)
func Bridge[T any](done <-chan interface{}, chanStream <-chan <-chan T) <-chan T {
	valStream := make(chan T)
	go func() {
		defer close(valStream)
		for {
			var stream <-chan T
			select {
			case mybeStream, ok := <-chanStream: //读取chanStream中的channel
				if !ok {
					return
				}
				stream = mybeStream
			case <-done:
				return
			}
			for val := range OrDone(done, stream) { //读取channel内容发送回去
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

// Tee 读取in数据并同时发送给两个接受的channel
func Tee[T any](done <-chan interface{}, in <-chan T) (_, _ <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)
	go func() {
		defer close(out1)
		defer close(out2)
		for v := range OrDone(done, in) {
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

// OrDone 通过done来控制性读取chan
func OrDone[T any](done <-chan interface{}, c <-chan T) <-chan T {
	valStream := make(chan T)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select { //可以进行优化
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

// MyInteger 整数类型，用于随机数类型转换
type MyInteger interface {
	~int | ~int32 | ~int64
}

// MyFloat 浮点数，可用于加减乘除
type MyFloat interface {
	~float64 | ~float32
}

// Number 可以用于加减乘除
type Number interface {
	MyFloat | MyInteger
}

// FanIn 从多个channels中合并数据到一个channel
func FanIn[T any](done <-chan interface{}, channels []<-chan T) <-chan T {
	var wg sync.WaitGroup
	multiplexedStream := make(chan T)
	multiplex := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()
	return multiplexedStream
}

// Multiply 乘法
func Multiply[V Number](done <-chan interface{}, valueStream <-chan V, multiplier V) <-chan V {
	results := make(chan V)
	go func() {
		defer close(results)
		for v := range valueStream {
			select {
			case <-done:
				return
			case results <- v * multiplier:
			}
		}
	}()
	return results
}

// Add 加法
func Add[V Number](done <-chan interface{}, valueStream <-chan V, additive V) <-chan V {
	results := make(chan V)
	go func() {
		defer close(results)
		for v := range valueStream {
			select {
			case <-done:
				return
			case results <- v + additive:
			}
		}
	}()
	return results
}

// ToType 显式转换为对应类型
func ToType[T MyInteger](done <-chan interface{}, valueStream <-chan interface{}) <-chan T {
	stringStream := make(chan T)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(T):
			}
		}
	}()
	return stringStream
}

// PrimeFinder 获取并判断素数
func PrimeFinder[T MyInteger](done <-chan interface{}, intStream <-chan T) <-chan interface{} {
	results := make(chan interface{})
	go func() {
		defer close(results)
		for v := range intStream {
			select {
			case <-done:
				return
			default:
			}
			for i := T(2); i*i < v; i++ {
				if v%i == 0 {
					goto next
				}
			}
			results <- v
		next:
		}
	}()
	return results
}

// Take 取出num个数后结束
func Take[T any](done <-chan interface{}, valueStream <-chan T, num int) <-chan T {
	results := make(chan T)
	go func() {
		defer close(results)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case results <- <-valueStream:
			}
		}
	}()
	return results
}

// RepeatFn 重复调用函数
func RepeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	results := make(chan interface{})
	go func() {
		defer close(results)
		for {
			select {
			case <-done:
				return
			case results <- fn():
			}
		}
	}()
	return results
}

// Repeat 重复生成值
func Repeat[T any](done <-chan interface{}, values ...T) <-chan T {
	valueStream := make(chan T)
	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}
