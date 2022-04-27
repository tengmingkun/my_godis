package lock

import (
	"sort"
	"sync"
)

type Locks struct {
	table []*sync.RWMutex
}

func NewLocks(size int) *Locks {
	table := make([]*sync.RWMutex, size)
	for i := 0; i < size; i++ {
		rwmutex := &sync.RWMutex{}
		table[i] = rwmutex
	}
	return &Locks{table: table}
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	for _, v := range key {
		hash = hash * uint32(16777619)
		hash = hash * uint32(v)
	}
	return hash
}

func (l *Locks) spread(hashCode uint32) uint32 {
	tablesize := uint32(len(l.table))
	return (tablesize - 1) & hashCode
}

func (l *Locks) Lock(key string) {
	index := l.spread(fnv32(key))
	mu := l.table[index]
	mu.Lock()
}

func (l *Locks) UnLock(key string) {
	index := l.spread(fnv32(key))
	mu := l.table[index]
	mu.Unlock()
}

func (l *Locks) RLock(key string) {
	index := l.spread(fnv32(key))
	l.table[index].RLock()
}

func (l *Locks) RUnLock(key string) {
	index := l.spread(fnv32(key))
	l.table[index].RUnlock()
}

func (l *Locks) toLockIndeies(key []string, reverse bool) []uint32 {
	indices := make([]uint32, len(key))
	for i := 0; i < len(key); i++ {
		indices[i] = l.spread(fnv32(key[i]))
	}
	sort.Slice(indices, func(i, j int) bool {
		if !reverse {
			return indices[i] < indices[j]
		} else {
			return indices[i] > indices[j]
		}
	})
	return indices
}

func (l *Locks) Locks(key ...string) {
	indices := l.toLockIndeies(key, false)
	for _, index := range indices {
		l.table[index].Lock()
	}
}

func (l *Locks) RLocks(key ...string) {
	indices := l.toLockIndeies(key, false)
	for _, index := range indices {
		l.table[index].RLock()
	}
}

func (l *Locks) ULocks(key ...string) {
	indices := l.toLockIndeies(key, false)
	for _, index := range indices {
		l.table[index].Unlock()
	}
}

func (l *Locks) RUnlocks(key ...string) {
	indices := l.toLockIndeies(key, false)
	for _, index := range indices {
		l.table[index].RUnlock()
	}
}
