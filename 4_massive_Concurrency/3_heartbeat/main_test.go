package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDoWork2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	intSlice := []int{0, 1, 2, 3, 4}
	heartBeat, results := DoWork2(ctx, intSlice...)
	<-heartBeat // 等待go程开始处理的信号
	i := 0
	for r := range results {
		if want := intSlice[i]; r != want {
			require.Equal(t, want, r, "idx=", i)
		}
		i++
	}
}

func TestDoWork3(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	intSlice := []int{0, 1, 2, 3, 4}
	const timeout = 2 * time.Second
	heartBeat, results := DoWork3(ctx, timeout, intSlice...)
	<-heartBeat
	i := 0
	for {
		select {
		case r, ok := <-results:
			if ok == false {
				return
			}
			require.Equal(t, intSlice[i], r)
			i++
		case <-heartBeat: // 接收心跳 防止超时
		case <-time.After(timeout):
			t.Fatal("time out")
		}
	}
}
