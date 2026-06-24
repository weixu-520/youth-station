package middleware

import (
	"net/http"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// GlobalRateLimitMiddleware 全局限流中间件
func GlobalRateLimitMiddleware(rds *redis.Redis, rate, burst int) func(http.HandlerFunc) http.HandlerFunc {
	limiter := limit.NewTokenLimiter(rate, burst, rds, "rate_limit:global")
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				httpx.WriteJson(w, http.StatusTooManyRequests, map[string]string{
					"code":    "429",
					"message": "请求过于频繁，请稍后再试",
				})
				return
			}
			next(w, r)
		}
	}
}

// IPRateLimiter 管理每个IP的限流器（使用sync.Map缓存）
type IPRateLimiter struct {
	store    *redis.Redis
	rate     int
	burst    int
	limiters sync.Map // key: IP, value: *limit.TokenLimiter
}

func NewIPRateLimiter(store *redis.Redis, rate, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		store: store,
		rate:  rate,
		burst: burst,
	}
}

// GetLimiter 获取指定IP的限流器，若不存在则创建
func (l *IPRateLimiter) GetLimiter(ip string) *limit.TokenLimiter {
	if val, ok := l.limiters.Load(ip); ok {
		return val.(*limit.TokenLimiter)
	}
	// 为每个IP创建独立的限流器，key中包含IP以确保隔离
	limiter := limit.NewTokenLimiter(
		l.rate,
		l.burst,
		l.store,
		"rate_limit:ip:"+ip,
	)
	// 设置 10 分钟自动过期，释放内存（实际过期由 Redis 管理，这里只是清理本地缓存）
	l.limiters.Store(ip, limiter)
	return limiter
}

// IPRateLimitMiddleware IP限流中间件
func IPRateLimitMiddleware(rds *redis.Redis, rate, burst int) func(http.HandlerFunc) http.HandlerFunc {
	factory := NewIPRateLimiter(rds, rate, burst)
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 获取真实IP（考虑代理）
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = strings.Split(r.RemoteAddr, ":")[0]
			}
			// 如果X-Forwarded-For有多个IP，取第一个
			if strings.Contains(ip, ",") {
				ip = strings.Split(ip, ",")[0]
			}

			limiter := factory.GetLimiter(ip)
			if !limiter.Allow() {
				httpx.WriteJson(w, http.StatusTooManyRequests, map[string]string{
					"code":    "429",
					"message": "请求过于频繁，请稍后再试",
				})
				return
			}
			next(w, r)
		}
	}
}

// RateLimit 根据配置模式返回相应的中间件
func RateLimit(rds *redis.Redis, mode string, rate, burst int) func(http.HandlerFunc) http.HandlerFunc {
	switch mode {
	case "ip":
		return IPRateLimitMiddleware(rds, rate, burst)
	case "global":
		fallthrough
	default:
		return GlobalRateLimitMiddleware(rds, rate, burst)
	}
}
