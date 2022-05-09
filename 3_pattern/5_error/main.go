package main

import (
	"fmt"
	"net/http"
)

/*
	处理并发程序下的错误
*/

type Result struct {
	Error    error
	Response *http.Response
}

func T1() {
	checkStatus := func(done <-chan interface{}, urls []string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				resp, err := http.Get(url)
				result := Result{Error: err, Response: resp}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}
	done := make(chan interface{})
	defer close(done)
	urls := []string{"https://google.com", "https://badhost"}
	for result := range checkStatus(done, urls) {
		if result.Error != nil {
			//可以加入一些对错误的判断，如对数量的判断
			fmt.Println("error: ", result.Error)
			continue
		}
		fmt.Println("result: ", result.Response)
	}
}

func main() {
	T1()
}
