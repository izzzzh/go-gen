package limit

import (
	"sync"
	"time"
)

// TokenBucket 表示一个令牌桶
type TokenBucket struct {
	mu         sync.Mutex
	capacity   int64     // 桶的容量
	tokens     int64     // 当前令牌数量
	lastRefill time.Time // 上次添加令牌的时间
	rate       float64   // 每秒添加的令牌数
}

// NewTokenBucket 创建一个新的令牌桶
func NewTokenBucket(capacity int64, rate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     0,
		lastRefill: time.Now(),
		rate:       rate,
	}
}

// Allow 尝试从桶中获取一个令牌，如果成功返回true，否则返回false
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	// 计算自上次添加令牌以来应该有多少新的令牌产生
	newTokens := int64(tb.rate * now.Sub(tb.lastRefill).Seconds())
	if newTokens > 0 {
		// 更新令牌数量，但不超过桶的容量
		tb.tokens = min(tb.tokens+newTokens, tb.capacity)
		tb.lastRefill = now
	}

	// 尝试从桶中获取一个令牌
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// min 返回两个整数中较小的一个
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	// 创建一个令牌桶，容量为10，每秒产生1个令牌
	tb := NewTokenBucket(10, 1.0)

	// 模拟请求并尝试获取令牌
	for i := 0; i < 20; i++ {
		if tb.Allow() {
			println("Request", i, "is allowed")
		} else {
			println("Request", i, "is rate limited")
		}
		time.Sleep(500 * time.Millisecond) // 暂停0.5秒以模拟请求间隔
	}
}
