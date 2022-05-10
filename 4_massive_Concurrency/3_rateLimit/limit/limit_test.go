package limit

import (
	"context"
	"strconv"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

const LAYOUT = "2006 01-02 15:04:05"

// 使用golang/rate 实现, 令牌桶算法
func TestRateLimitByGoRate1(t *testing.T) {
	ticker := rate.NewLimiter(3, 6)
	length := 20
	chs := make([]chan string, length)
	for i := 0; i < length; i++ {
		chs[i] = make(chan string, 1)
		go func(taskId string, ch chan string, r *rate.Limiter) {
			err := r.Wait(context.Background())
			if err != nil {
				ch <- "Task-" + taskId + " not allow " + time.Now().Format(LAYOUT)
			}

			time.Sleep(time.Duration(5) * time.Millisecond)
			ch <- "Task-" + taskId + " run success  " + time.Now().Format(LAYOUT)
			return

		}(strconv.FormatInt(int64(i), 10), chs[i], ticker)
	}
	for _, ch := range chs {
		t.Log("task start at " + <-ch)
	}
}
