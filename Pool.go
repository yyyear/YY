package YY

import "sync"

type Pool[T PoolItem[T]] struct {
	p sync.Pool
}

// PoolItem 需要实现接口
type PoolItem[T any] interface {
	Reset()
}

func NewPool[T PoolItem[T]](t T) *Pool[T] {
	return &Pool[T]{
		p: sync.Pool{New: func() any {
			return t
		}},
	}
}

func (p *Pool[T]) Get() T {
	return p.p.Get().(T)
}

func (p *Pool[T]) Put(t T) {
	t.Reset()
	p.p.Put(t)
}
