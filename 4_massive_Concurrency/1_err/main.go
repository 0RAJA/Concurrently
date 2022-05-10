package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/0RAJA/Concurrently/4_massive_Concurrency/1_err/myerr"
)

func handleErr(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID:%v]:", key))
	log.Printf("%#v", err)
	fmt.Printf("[%v] %v", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.Llongfile)
	err := myerr.RunJob("1")
	if err != nil {
		msg := "there was an unexpected issue"
		if errors.As(err, &myerr.IntermediateErr{}) { // 如果这个是预期的包装良好的错误，则可以直接传递给用户
			msg = err.Error()
		}
		handleErr(1, err, msg)
	}
}
