package set

import "test/datastruct/dict"

type Set struct {
	dict dict.Dict
}

func NewSet() *Set {
	set := dict.NewSimpleDict()
	return &Set{dict: set}
}

func NewSetWithVal(val ...string) *Set {
	set := &Set{dict: dict.NewConcurrentDict(len(val))}
	for _, v := range val {
		set.Add(v)
	}
	return set
}

func (s *Set) Add(key string) int {
	return s.dict.Put(key, true)
}

func (s *Set) Remove(key string) int {
	return s.dict.Remove(key)
}
func (s *Set) Has(key string) bool {
	_, statu := s.dict.Get(key)
	return statu
}
func (s *Set) Len() int {
	return s.dict.Len()
}
func (s *Set) ToSlice() []string {
	return s.dict.Keys()
}
func (s *Set) ForEach(consumer func(menber string) bool) {
	s.dict.ForEach(func(key string, val interface{}) bool {
		return consumer(key)
	})

}
func (set *Set) Intersect(another *Set) *Set {
	if set == nil {
		panic("set is nil")
	}

	result := NewSet()
	another.ForEach(func(member string) bool {
		if set.Has(member) {
			result.Add(member)
		}
		return true
	})
	return result
}
func (set *Set) Union(another *Set) *Set {
	if set == nil {
		panic("set is nil")
	}
	result := NewSet()
	another.ForEach(func(member string) bool {
		result.Add(member)
		return true
	})
	set.ForEach(func(member string) bool {
		result.Add(member)
		return true
	})
	return result
}

func (set *Set) Diff(another *Set) *Set {
	if set == nil {
		panic("set is nil")
	}

	result := NewSet()
	set.ForEach(func(member string) bool {
		if !another.Has(member) {
			result.Add(member)
		}
		return true
	})
	return result
}

func (set *Set) RandomMembers(limit int) []string {
	return set.dict.RandomKeys(limit)
}

func (set *Set) RandomDistinctMembers(limit int) []string {
	return set.dict.RandomDistinctKeys(limit)
}
