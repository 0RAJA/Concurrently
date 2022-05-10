package main

import (
	"context"
	"log"
	"os"
	"sync"

	"Concurrently/4_massive_Concurrency/3_rateLimit/api"
)

/*
	速率限制
*/

func main() {
	defer log.Println("Done")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.Lshortfile)
	apiConn := api.Open()
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConn.ReadFile(context.Background()); err != nil {
				log.Println("cannot read file:", err)
				return
			}
			log.Println("read file")
		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConn.ResolveAddress(context.Background()); err != nil {
				log.Println("cannot resolve address:", err)
				return
			}
			log.Println("ResolveAddress")
		}()
	}
	wg.Wait()
}
