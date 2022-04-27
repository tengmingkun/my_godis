package dict

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
)

type ConcurrentDict struct {
	table []*Shared
	count int32
}
type Shared struct {
	m  map[string]interface{}
	mu sync.RWMutex
}

func computeCapacity(param int) (size int) {
	if param <= 16 {
		return 16
	}
	n := param - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return math.MaxInt32
	} else {
		return int(n + 1)
	}
}
func NewConcurrentDict(sharesCount int) *ConcurrentDict {
	count := computeCapacity(sharesCount)
	//fmt.Println(count)
	table := make([]*Shared, count)
	for i := 0; i < count; i++ {
		table[i] = &Shared{m: make(map[string]interface{})}
	}
	return &ConcurrentDict{
		table: table,
		count: 0,
	}
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	for _, v := range key {
		hash = hash * uint32(16777619)
		hash = hash * uint32(v)
	}
	return hash
}

func (d *ConcurrentDict) spread(hashCode uint32) uint32 {
	tablesize := len(d.table)
	return uint32(tablesize-1) & hashCode
}

func (d *ConcurrentDict) getShared(index uint32) *Shared {
	return d.table[index]
}

func (d *ConcurrentDict) Get(key string) (interface{}, bool) {
	shared := d.getShared(d.spread(fnv32(key)))
	shared.mu.RLock()
	defer shared.mu.RUnlock()
	if val, ok := shared.m[key]; ok {
		return val, ok
	} else {
		return nil, false
	}
}

func (d *ConcurrentDict) Len() int {
	return int(atomic.LoadInt32(&d.count))
}

func (d *ConcurrentDict) Put(key string, val interface{}) int {
	shared := d.getShared(d.spread(fnv32(key)))
	if shared == nil {
		fmt.Println("shared is nil")
	}
	shared.mu.Lock()
	defer shared.mu.Unlock()
	if _, ok := shared.m[key]; ok {
		shared.m[key] = val
		return 0
	} else {
		shared.m[key] = val
		d.AddCount()
		return 1
	}
}
func (d *ConcurrentDict) AddCount() {
	atomic.AddInt32(&d.count, 1)
}
func (d *ConcurrentDict) PutIfAbsent(key string, val interface{}) int {
	shared := d.getShared(d.spread(fnv32(key)))
	shared.mu.Lock()
	defer shared.mu.Unlock()
	if _, ok := shared.m[key]; ok {
		return 0
	} else {
		shared.m[key] = val
		d.AddCount()
		return 1
	}
}

func (d *ConcurrentDict) PutIfExists(key string, val interface{}) int {
	shared := d.getShared(d.spread(fnv32(key)))
	shared.mu.Lock()
	defer shared.mu.Unlock()
	if _, ok := shared.m[key]; ok {
		shared.m[key] = val
		return 1
	} else {
		return 0
	}
}

func (d *ConcurrentDict) Remove(key string) int {
	shared := d.getShared(d.spread(fnv32(key)))
	shared.mu.Lock()
	defer shared.mu.Unlock()
	if _, ok := shared.m[key]; ok {
		delete(shared.m, key)
		atomic.AddInt32(&d.count, -1)
		return 1
	} else {
		return 0
	}
}

func (d *ConcurrentDict) ForEach(consumer Consumer) {
	for _, shared := range d.table {
		for key, val := range shared.m {
			shared.mu.Lock()
			continues := consumer(key, val)
			shared.mu.Unlock()
			if !continues {
				return
			}
		}
	}
}

func (d *ConcurrentDict) Keys() []string {
	//fmt.Println("应该的长度", d.Len())
	keys := make([]string, d.Len())
	i := 0
	d.ForEach(func(key string, val interface{}) bool {
		if i < len(keys) {
			keys[i] = key
			i++
		} else {
			keys = append(keys, key)
		}
		return true
	})
	return keys
}

func (shard *Shared) RandomKey() string {
	if shard == nil {
		panic("shard is nil")
	}
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	for key := range shard.m {
		return key
	}
	return ""
}

func (dict *ConcurrentDict) RandomKeys(limit int) []string {
	size := dict.Len()
	if limit >= size {
		return dict.Keys()
	}
	shardCount := len(dict.table)

	result := make([]string, limit)
	for i := 0; i < limit; {
		shard := dict.getShared(uint32(rand.Intn(shardCount)))
		if shard == nil {
			continue
		}
		key := shard.RandomKey()
		if key != "" {
			result[i] = key
			i++
		}
	}
	return result
}

func (dict *ConcurrentDict) RandomDistinctKeys(limit int) []string {
	size := dict.Len()
	if limit >= size {
		return dict.Keys()
	}

	shardCount := len(dict.table)
	result := make(map[string]bool)
	for len(result) < limit {
		shardIndex := uint32(rand.Intn(shardCount))
		shard := dict.getShared(shardIndex)
		if shard == nil {
			continue
		}
		key := shard.RandomKey()
		if key != "" {
			result[key] = true
		}
	}
	arr := make([]string, limit)
	i := 0
	for k := range result {
		arr[i] = k
		i++
	}
	return arr
}
