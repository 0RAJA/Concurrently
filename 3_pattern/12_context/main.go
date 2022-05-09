package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

/*
	使用context存储value的要求：
	仅将上下文值用于传输进程和请求的请求范围数据，API边界，而不是将可选参数传递给函数
		启发式：
	1. 数据应该通过进程或者API边界
		如果在进程的内存中生成数据，除非是通过API边界传递数据，否则不是很好。
	2. 数据应该是不变的
	3. 数据应该趋向于简单类型
		便于使用方拉出数据
	4. 数据应该是数据，而不是方法或类型
	5. 数据应该有助于修饰操作，而不是驱动它们
		算法不应该依赖于context中的数据。
*/

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}

func main() {
	t2()
}

func t2() {
	ProcessRequest("1", "2")
}

func ProcessRequest(userID, auth string) {
	ctx := context.WithValue(context.Background(), "userID", userID)
	ctx = context.WithValue(ctx, "auth", auth)
	handleResponse(ctx)
}

func handleResponse(ctx context.Context) {
	fmt.Println(ctx.Value("userID"), ctx.Value("auth"))
}

//一个打招呼的函数
func t1() {
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(ctx); err != nil {
			log.Println(err)
			cancel()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(ctx); err != nil {
			log.Println(err)
		}
	}()
	wg.Wait()
}
func printFarewell(ctx context.Context) error {
	farewell, err := genFarewell(ctx)
	if err != nil {
		return err
	}
	fmt.Println(farewell)
	return nil
}
func printGreeting(ctx context.Context) error {
	greeting, err := genGreeting(ctx)
	if err != nil {
		return err
	}
	fmt.Println(greeting)
	return nil
}
func genGreeting(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second) //做一个超时判断
	defer cancel()
	switch l, err := locale(ctx); {
	case err != nil:
		return "", err
	case l == "EN/US":
		return "hello", nil
	}
	return "", errors.New("unsupported")
}

func genFarewell(ctx context.Context) (string, error) {
	switch l, err := locale(ctx); {
	case err != nil:
		return "", err
	case l == "EN/US":
		return "goodbye", nil
	}
	return "", errors.New("unsupported")
}
func locale(ctx context.Context) (string, error) {
	if deadline, ok := ctx.Deadline(); ok { //提前检查时间
		if deadline.Before(time.Now().Add(2 * time.Second)) {
			return "", context.DeadlineExceeded
		}
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(2 * time.Second):
		return "EN/US", nil
	}
}
