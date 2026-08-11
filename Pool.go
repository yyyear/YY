package YY

import "sync"

type PoolItem interface {
	Reset()
}

type Pool[T PoolItem] struct {
	pool sync.Pool
}

func NewPool[T PoolItem](newItem func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any { return newItem() },
		},
	}
}

func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

func (p *Pool[T]) Put(item T) {
	item.Reset()
	p.pool.Put(item)
}
