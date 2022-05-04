package zset

import (
	"fmt"
	"test/datastruct/skiplist"
)

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

func (zset *Zset) GetCount(min, max int) int {
	return zset.skl.GetCount(min, max)
}

func (zset *Zset) FindMember(val interface{}) (int, bool) {
	bytes, ok := zset.skl.FindMember(val)
	return bytes, ok
}

func (zset *Zset) GetRank(val interface{}) int {
	return zset.skl.GetRank(val)
}

func (zset *Zset) DeleteKey(start, end int) int {
	count := 0
	fmt.Println(start, end)
	for i := start; i <= end; i++ {
		node := zset.skl.Search(i)
		if node != nil {
			zset.skl.Delete(i)
			zset.Size--
			count++
		}
	}
	return count
}
