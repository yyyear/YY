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

type NSString string

// Builder 增加字符
func Builder() *strings.Builder {
	return builderPool.Get().(*strings.Builder)
}

func BuilderRelease(b *strings.Builder) {
	b.Reset()
	builderPool.Put(b)
}
