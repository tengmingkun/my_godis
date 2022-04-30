package zset

import "test/datastruct/skiplist"

type Zset struct {
	skl  *skiplist.SkipList
	Size int
}

func NewZset() *Zset {
	return &Zset{skl: skiplist.NewSkipList()}
}

func (zset *Zset) ADD(socre int, val interface{}) bool {
	if zset == nil {
		return false
	}
	zset.skl.Insert(socre, val)
	zset.Size++
	return true
}

func (zset *Zset) GetRange(start int, end int) ([]int, []interface{}) {
	if zset == nil {
		return nil, nil
	}
	if start < 0 {
		start = 0
	}
	if end >= zset.Size {
		end = zset.Size - 1
	}
	return zset.skl.GetRange(start, end)
}
