package dict

import "fmt"

type SimpleDict struct {
	m map[string]interface{}
}

func NewSimpleDict() *SimpleDict {
	m := make(map[string]interface{})
	return &SimpleDict{m: m}
}

func (dict *SimpleDict) Get(key string) (val interface{}, ok bool) {
	val, ok = dict.m[key]
	return
}
func (dict *SimpleDict) Len() int {
	return len(dict.m)
}

func (dict *SimpleDict) Put(key string, val interface{}) int {
	_, existed := dict.m[key]
	dict.m[key] = val
	if existed {
		return 0
	} else {
		return 1
	}
}
func (dict *SimpleDict) PutIfAbsent(key string, val interface{}) int {
	_, existed := dict.m[key]
	if existed {
		return 0
	} else {
		dict.m[key] = val
		return 1
	}
}

func (dict *SimpleDict) PutIfExists(key string, val interface{}) int {
	_, existed := dict.m[key]
	if !existed {
		return 0
	} else {
		dict.m[key] = val
		return 1
	}
}

func (dict *SimpleDict) Remove(key string) int {
	_, existed := dict.m[key]
	if !existed {
		return 0
	} else {
		delete(dict.m, key)
		return 1
	}
}

func (dict *SimpleDict) Keys() []string {
	//fmt.Println("应该的长度", dict.Len())
	result := make([]string, 0)
	for k := range dict.m {
		result = append(result, k)
	}
	//fmt.Println(result, len(result))
	return result
}

func (dict *SimpleDict) ForEach(consumer Consumer) {
	for k, v := range dict.m {
		if !consumer(k, v) {
			break
		}
	}
}
func (dict *SimpleDict) RandomKeys(limit int) []string {
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		for k := range dict.m { //range 的每次遍历都是不一样的
			result[i] = k
			break
		}
	}
	fmt.Println(result)
	return result
}

func (dict *SimpleDict) RandomDistinctKeys(limit int) []string {
	size := limit
	if size > len(dict.m) {
		size = len(dict.m)
	}
	result := make([]string, size)
	i := 0
	for k := range dict.m {
		if i == limit {
			break
		}
		result[i] = k
		i++
	}
	return result
}
