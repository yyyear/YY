package YY

import (
	"strings"
	"sync"
)

var builderPool sync.Pool

func init() {
	builderPool = sync.Pool{New: func() interface{} {
		return Build{}
	}}
}

type Build strings.Builder

func NewBuilder() Build {
	return builderPool.Get().(Build)
}

// Add 增加字符
func (b Build) Add(r string) Build {
	var j strings.Builder = strings.Builder(b)
	j.WriteString(r)
	return Build(j)
}
func (b Build) String() string {
	builder := strings.Builder(b)
	return builder.String()
}
func (b Build) Release() {
	builder := strings.Builder(b)
	builder.Reset()
	builderPool.Put(builder)
}

func BuilderRelease(b *strings.Builder) {
	b.Reset()
	builderPool.Put(b)
}
