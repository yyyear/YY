package YY

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Value map[string]any

// ---------- 基础方法（可以用接收器，因为没有类型参数） ----------

func (v Value) Set(key string, val any) Value {
	v[key] = val
	return v
}

func (v Value) Has(key string) bool {
	_, ok := v[key]
	return ok
}

func (v Value) Delete(key string) {
	delete(v, key)
}
func (v Value) String(key string) string {
	vvv, ok := GetOK[string](v, key)
	if ok {
		return vvv
	}
	return ""
}
func (v Value) Int(key string) int {
	vvv, ok := GetOK[int](v, key)
	if ok {
		return vvv
	}
	return 0
}
func (v Value) Int32(key string) int32 {
	vvv, ok := GetOK[int32](v, key)
	if ok {
		return vvv
	}
	return 0
}
func (v Value) Data() []byte {
	data, _ := json.Marshal(v)
	return data
}

// GetOK 原来带 (T, bool) 的版本
func GetOK[T any](v Value, key string) (T, bool) {
	var zero T
	raw, ok := v[key]
	if !ok {
		return zero, false
	}
	return cast[T](raw)
}

// ---------- 泛型取值函数 ----------

// Get 按 key 取值并断言为 T；不存在或类型不匹配时返回零值和 false
func Get[T any](v Value, key string) (T, bool) {
	var zero T
	raw, ok := v[key]
	if !ok {
		return zero, false
	}
	return cast[T](raw)
}

// GetOr 按 key 取值，失败返回默认值
func GetOr[T any](v Value, key string, def T) T {
	if t, ok := Get[T](v, key); ok {
		return t
	}
	return def
}

// MustGet 按 key 取值，失败直接 panic
func MustGet[T any](v Value, key string) T {
	t, ok := Get[T](v, key)
	if !ok {
		panic(fmt.Sprintf("value: key %q not found or type mismatch", key))
	}
	return t
}

// GetE 按 key 取值，失败返回错误
func GetE[T any](v Value, key string) (T, error) {
	var zero T
	raw, ok := v[key]
	if !ok {
		return zero, fmt.Errorf("value: key %q not found", key)
	}
	t, ok := cast[T](raw)
	if !ok {
		return zero, fmt.Errorf("value: key %q type mismatch: got %T, want %T", key, raw, zero)
	}
	return t, nil
}

// Path 支持 "a.b.c" 形式的嵌套取值
func Path[T any](v Value, path string) (T, bool) {
	var zero T
	var cur any = map[string]any(v)
	for _, p := range strings.Split(path, ".") {
		m, ok := toMap(cur)
		if !ok {
			return zero, false
		}
		cur, ok = m[p]
		if !ok {
			return zero, false
		}
	}
	return cast[T](cur)
}

// ---------- 内部辅助 ----------

func toMap(a any) (map[string]any, bool) {
	switch m := a.(type) {
	case map[string]any:
		return m, true
	case Value:
		return map[string]any(m), true
	}
	return nil, false
}

// cast 做类型断言，并对常见场景做兼容：
//   - JSON 反序列化后数字都是 float64，可自动转 int / int64 / float64
//   - json.Number / 字符串形式的数字也可转换
//   - []byte 转 string
func cast[T any](raw any) (T, bool) {
	var zero T
	if t, ok := raw.(T); ok {
		return t, true
	}

	switch any(zero).(type) {
	case int:
		switch n := raw.(type) {
		case float64:
			return any(int(n)).(T), true
		case int64:
			return any(int(n)).(T), true
		case json.Number:
			if i, err := n.Int64(); err == nil {
				return any(int(i)).(T), true
			}
		case string:
			if i, err := strconv.Atoi(n); err == nil {
				return any(i).(T), true
			}
		}
	case int64:
		switch n := raw.(type) {
		case float64:
			return any(int64(n)).(T), true
		case int:
			return any(int64(n)).(T), true
		case json.Number:
			if i, err := n.Int64(); err == nil {
				return any(i).(T), true
			}
		case string:
			if i, err := strconv.ParseInt(n, 10, 64); err == nil {
				return any(i).(T), true
			}
		}
	case float64:
		switch n := raw.(type) {
		case int:
			return any(float64(n)).(T), true
		case int64:
			return any(float64(n)).(T), true
		case json.Number:
			if f, err := n.Float64(); err == nil {
				return any(f).(T), true
			}
		case string:
			if f, err := strconv.ParseFloat(n, 64); err == nil {
				return any(f).(T), true
			}
		}
	case string:
		switch n := raw.(type) {
		case []byte:
			return any(string(n)).(T), true
		case fmt.Stringer:
			return any(n.String()).(T), true
		}
	case bool:
		if s, ok := raw.(string); ok {
			if b, err := strconv.ParseBool(s); err == nil {
				return any(b).(T), true
			}
		}
	}

	return zero, false
}
