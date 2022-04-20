package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

/*
	Pool模式是一种创建和提供可供使用的固定数量实例或Pool实例的方法。
	它常用于约束创建昂贵的场景下。
	通过Get()来获取新的实例，通过Put()释放资源
	可以用于快速将预先分配的对象缓存加载启动。
*/

func T1() {
	myPool := &sync.Pool{New: func() interface{} {
		fmt.Println("create object")
		return struct{}{}
	}}
	a := myPool.Get()
	b := myPool.Get()
	myPool.Put(a)
	myPool.Put(b)
	myPool.Get()
}

func connectHandle() interface{} {
	time.Sleep(time.Second)
	return struct{}{}
}

func warmServiceConnCache() *sync.Pool {
	p := &sync.Pool{
		New: connectHandle,
	}
	for i := 0; i < 10; i++ {
		p.Put(p.New())
	}
	return p
}

func startNetWorkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		connPool := warmServiceConnCache()
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatal(err)
		}
		defer server.Close()
		wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Println("Cannot accept connection")
				continue
			}
			go func(conn net.Conn) {
				defer conn.Close()
				svcConn := connPool.Get()
				fmt.Fprintln(conn, "")
				connPool.Put(svcConn)
			}(conn)
		}
	}()
	return &wg
}

func T2() {

}

func main() {
	//T1()
}
