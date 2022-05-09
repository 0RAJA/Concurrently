package main

/*
	处理来自系统各个部分的channel,你不知道你的go程是否被取消。
*/
func OrDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
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

func main() {

}
