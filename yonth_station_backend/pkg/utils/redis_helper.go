package utils

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

func CacheAside(ctx context.Context, rds *redis.Redis, key string, expire int, queryFunc func() (interface{}, error)) (interface{}, error) {
	// 1. 查缓存
	val, err := rds.GetCtx(ctx, key)
	if err == nil && val != "" {
		return val, nil // 返回字符串，由调用方自行反序列化
	}
	// 2. 查 DB
	data, err := queryFunc()
	if err != nil {
		return nil, err
	}
	// 3. 写入缓存（序列化由调用方处理）
	// 此处仅返回 data，由调用方决定是否序列化存入
	return data, nil
}

// DeleteCache 删除缓存（用于更新后失效）
func DeleteCache(ctx context.Context, rds *redis.Redis, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	_, err := rds.DelCtx(ctx, keys...)
	return err
}

func DecrWithMin(ctx context.Context, rds *redis.Redis, key string, min int64) (int64, error) {
	// 使用 Lua 脚本保证原子性
	script := `
        local current = redis.call('GET', KEYS[1])
        if not current then
            return -1  -- key不存在
        end
        local newVal = tonumber(current) - 1
        if newVal < tonumber(ARGV[1]) then
            return -2  -- 低于最小值
        end
        redis.call('SET', KEYS[1], newVal)
        return newVal
    `
	result, err := rds.EvalCtx(ctx, script, []string{key}, min)
	if err != nil {
		return 0, err
	}
	val, ok := result.(int64)
	if !ok {
		return 0, errors.New("unexpected result type")
	}
	if val == -1 {
		return 0, errors.New("key not exists")
	}
	if val == -2 {
		return 0, errors.New("insufficient quota")
	}
	return val, nil
}

// IncrBy 原子增加
func IncrBy(ctx context.Context, rds *redis.Redis, key string, delta int64) (int64, error) {
	return rds.IncrbyCtx(ctx, key, delta)
}

// SetEx 封装
func SetEx(ctx context.Context, rds *redis.Redis, key string, value interface{}, expire int) error {
	val, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	return rds.SetexCtx(ctx, key, val, expire)
}

// GetString 获取字符串值
func GetString(ctx context.Context, rds *redis.Redis, key string) (string, error) {
	return rds.GetCtx(ctx, key)
}
