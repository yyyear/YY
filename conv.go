package YY

import (
	"fmt"
	"strconv"
)

// ToInt64Base 泛型转换字符到数字
func ToInt64Base(s string, base int) int64 {
	i, err := strconv.ParseInt(s, base, 64)
	if err != nil {
		Warn("ParseInt64 failed", Fields{"input": s, "error": err.Error()})
	}
	return i
}
func ToInt64(s string) int64 {
	return ToInt64Base(s, 10)
}

// ToUInt64 泛型转换字符到数字
func ToUInt64(s string) uint64 {
	return ToUInt64Base(s, 10)
}

// ToUInt64Base 泛型转换字符到数字
func ToUInt64Base(s string, base int) uint64 {
	i, err := strconv.ParseUint(s, base, 64)
	if err != nil {
		Warn("ParseInt64 failed", Fields{"input": s, "error": err.Error()})
	}
	return i
}

// Integer 包含有符号和无符号整数（可根据需要增减）
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// ToStr 泛型转换数字到字符
func ToStr[T Integer](num T) string {
	return ToStringBase(num, 10)
}

// ToStringBase 泛型转换数字到字符
func ToStringBase[T Integer](num T, base int) string {
	switch any(num).(type) {
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(int64(num), base)
	case uint, uint8, uint16, uint32, uint64, uintptr:
		return strconv.FormatUint(uint64(num), base)
	default:
		return fmt.Sprint(num) // 保底
	}
}
