package YY

import (
	"strings"
	"sync"
)

var builderPool sync.Pool

func init() {
	builderPool = sync.Pool{New: func() interface{} {
		return &strings.Builder{}
	}}
}

// Add 增加字符
func Add(l string, r string) string {
	builder := builderPool.Get().(*strings.Builder)
	builder.WriteString(l)
	builder.WriteString(r)
	result := builder.String()
	builder.Reset()
	builderPool.Put(builder)
	return result
}
