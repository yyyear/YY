package YY

import "hash/fnv"

// ShardedMap 分片哈希表，降低扩容成本与锁竞争（如果是只读，甚至不需要锁）
type ShardedMap[T any] struct {
	shardCount uint64
	// 使用固定大小的 struct，Value 只是切片的索引 (uint32)
	// 这个 Map 底层没有任何指针，GC 会直接跳过扫描！
	maps []map[uint64]uint32
	data []T // 真正的数据存储在连续的切片中
}

// NewShardedMap 创建分片哈希表
func NewShardedMap[T any](shards uint64, capacity int) *ShardedMap[T] {
	sm := &ShardedMap[T]{
		shardCount: shards,
		maps:       make([]map[uint64]uint32, shards),
		data:       make([]T, 0, capacity),
	}
	for i := range shards {
		sm.maps[i] = make(map[uint64]uint32, capacity/int(shards))
	}
	return sm
}

var fnvHash = fnv.New64a()

// getHash 将 string 转换为 uint64
func (sm *ShardedMap[T]) getHash(key string) uint64 {
	fnvHash.Reset()
	_, _ = fnvHash.Write([]byte(key))
	return fnvHash.Sum64()
}

// Put 插入数据
func (sm *ShardedMap[T]) Put(key string, p T) {
	hash := sm.getHash(key)
	shard := hash % sm.shardCount

	// 将数据追加到连续切片中，拿到索引
	sm.data = append(sm.data, p)
	idx := uint32(len(sm.data) - 1)

	sm.maps[shard][hash] = idx
}

// Get 查找数据
func (sm *ShardedMap[T]) Get(key string) (T, bool) {
	hash := sm.getHash(key)
	shard := hash % sm.shardCount

	idx, exists := sm.maps[shard][hash]
	if !exists {
		return *new(T), false
	}
	return sm.data[idx], true
}
