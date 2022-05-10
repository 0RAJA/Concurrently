package api

import (
	"context"
	"time"

	"Concurrently/4_massive_Concurrency/3_rateLimit/limit"
	"golang.org/x/time/rate"
)

type API interface {
	ReadFile(ctx context.Context) error
	ResolveAddress(ctx context.Context) error
}

type testApi struct {
	netWorkLimit, diskLimit, apiLimit limit.RateLimiter // 多个维度进行限制
}

func Open() API {
	apiLimit := limit.MultiLimiter(
		rate.NewLimiter(limit.Per(2, time.Second), 1),   // 每秒的限制,防止突发请求,每1秒补充两个
		rate.NewLimiter(limit.Per(10, time.Minute), 10), // 每分钟的限制，设置初始池,每10秒补充一个
	)
	diskLimit := limit.MultiLimiter(
		rate.NewLimiter(rate.Limit(1), 1),
	)
	netWorkLimit := limit.MultiLimiter(
		rate.NewLimiter(limit.Per(3, time.Second), 3),
	)
	return &testApi{
		apiLimit:     apiLimit,
		diskLimit:    diskLimit,
		netWorkLimit: netWorkLimit,
	}
}

func (t *testApi) ReadFile(ctx context.Context) error {
	if err := limit.MultiLimiter(t.apiLimit, t.diskLimit).Wait(ctx); err != nil { // 融合api限流和磁盘限流
		return err
	}
	return nil
}

func (t *testApi) ResolveAddress(ctx context.Context) error {
	if err := limit.MultiLimiter(t.apiLimit, t.netWorkLimit).Wait(ctx); err != nil {
		return err
	}
	return nil
}
