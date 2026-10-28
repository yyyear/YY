package YY

import (
	"encoding/binary"
	"errors"
	"strconv"

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
		Error("Result错误处理:", r.Error)
		return *new(T)
	}
	return r.Value
}

// Encode proto Model 转 []byte
func Encode[T proto.Message](m T) ([]byte, error) {
	b, err := proto.Marshal(m)
	if err != nil {
		Error("转换数据错误:", err)
		return nil, err
	}
	return b, nil
}

// Decode 转换 proto Model
func Decode[T proto.Message](bytes []byte, m T) Result[T] {
	err := proto.Unmarshal(bytes, m)
	if err != nil {
		return Try(m, err)
	}
	return Try(m, nil)
}

// EncodeListPack proto Model 数组 打包成单个 []byte
// 格式: [uvarint len][payload][uvarint len][payload]...
func EncodeListPack[T proto.Message](ms []T) ([]byte, error) {
	var buf []byte
	var lenBuf [binary.MaxVarintLen64]byte
	for i, m := range ms {
		b, err := Encode(m)
		if err != nil {
			Error("转换数据错误["+strconv.Itoa(i)+"]:", err)
			return nil, err
		}
		n := binary.PutUvarint(lenBuf[:], uint64(len(b)))
		buf = append(buf, lenBuf[:n]...)
		buf = append(buf, b...)
	}
	return buf, nil
}

// DecodeListPack 单个 []byte 解包成 proto Model 数组
func DecodeListPack[T proto.Message](data []byte, newFn func() T) Result[[]T] {
	var out []T
	for len(data) > 0 {
		l, n := binary.Uvarint(data)
		if n <= 0 {
			return Try[[]T](nil, errors.New("invalid varint length"))
		}
		data = data[n:]
		if uint64(len(data)) < l {
			return Try[[]T](nil, errors.New("unexpected EOF"))
		}
		m := newFn()
		if err := proto.Unmarshal(data[:l], m); err != nil {
			return Try[[]T](nil, err)
		}
		out = append(out, m)
		data = data[l:]
	}
	return Try(out, nil)
}
