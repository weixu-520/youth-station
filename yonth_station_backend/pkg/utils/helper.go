package utils

//辅助函数
import (
	"context"
	"encoding/json"
	"errors"
)

// getUserId 从 context 中获取当前登录用户 ID（兼容 go-zero JWT 中间件存入的 json.Number 类型）
func GetUserId(ctx context.Context) (int64, error) {
	val := ctx.Value("userId")
	if val == nil {
		return 0, errors.New("userId not found")
	}
	num, ok := val.(json.Number)
	if !ok {
		return 0, errors.New("userId type error")
	}
	return num.Int64()
}

// getStatusDesc 根据申请状态码返回中文描述
func GetStatusDesc(status int8) string {
	switch status {
	case 0:
		return "待审核"
	case 1:
		return "已通过"
	case 2:
		return "已拒绝"
	case 3:
		return "已取消"
	case 4:
		return "已入住"
	case 5:
		return "已退房"
	default:
		return "未知"
	}
}

// maskPhone 手机号脱敏：保留前3后4，中间4位星号
func MaskPhone(phone *string) string {
	if phone == nil || *phone == "" {
		return ""
	}
	p := *phone
	if len(p) == 11 {
		return p[:3] + "****" + p[7:]
	}
	return p // 非标准长度不脱敏
}

// maskIdCard 身份证号脱敏：保留前6后4，中间8位星号
func MaskIdCard(idCard string) string {
	if len(idCard) == 18 {
		return idCard[:6] + "********" + idCard[14:]
	}
	if idCard != "" {
		return "****"
	}
	return ""
}

// getIsAdmin 从 context 中获取 isAdmin 字段（由 JWT 中间件注入）
func GetIsAdmin(ctx context.Context) (bool, error) {
	val := ctx.Value("isAdmin")
	if val == nil {
		return false, errors.New("isAdmin not found")
	}
	// 注意：go-zero 默认 JWT 中间件会将 claims 中的值存储为 interface{}，
	// 实际类型可能是 bool、float64 或 string，需要兼容处理
	switch v := val.(type) {
	case bool:
		return v, nil
	case float64:
		return v != 0, nil
	case string:
		return v == "true", nil
	default:
		return false, errors.New("invalid isAdmin type")
	}
}
