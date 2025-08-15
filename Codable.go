package YY

import (
	"google.golang.org/protobuf/proto"
)

type Result[T any] struct {
	Value T
	Error error
}

func Try[T any](value T, err error) Result[T] {
	return Result[T]{Value: value, Error: err}
}
func TryError(err error) Result[bool] {
	return Result[bool]{Value: err == nil, Error: err}
}

func (r Result[T]) Do() T {
	if r.Error != nil {
		ErrorLeve(1, "错误处理:", r.Error)
	}
	return r.Value
}

// Encode proto Model 转 []byte
func Encode[T proto.Message](m T) []byte {
	b, err := proto.Marshal(m)
	if err != nil {
		ErrorLeve(1, "转换数据错误:", err)
	}
	return b
}

// Decode 转换 proto Model
func Decode[T proto.Message](bytes []byte, m T) Result[T] {
	err := proto.Unmarshal(bytes, m)
	if err != nil {
		return Try(m, err)
	}
	return Try(m, nil)
}
