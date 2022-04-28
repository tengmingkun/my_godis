package db

import (
	"strconv"
	"test/datastruct/set"
	"test/interface/redis"
	"test/redis/reply"
)

func SADD(db *DB, args [][]byte) redis.Reply {
	if len(args) <= 2 {
		return reply.MakeArgNumErrReply("sadd")
	}
	key := args[1]
	Set, replys := db.getOrInfoSet(string(key))
	if replys != nil {
		return replys
	}
	for i := 2; i < len(args); i++ {
		Set.Add(string(args[i]))
	}
	return reply.MakeIntReply(int64(Set.Len()))
}

func SPOP(db *DB, args [][]byte) redis.Reply {
	if len(args) < 2 {
		return reply.MakeArgNumErrReply("spop")
	}
	key := args[1]
	Set, replys := db.getAsSet(string(key))
	if replys != nil {
		return replys
	}
	if Set == nil || Set.Len() == 0 {
		return reply.MakeIntReply(0)
	}
	num := 0
	var err error
	if len(args) == 2 {
		num = 1
	} else {
		num, err = strconv.Atoi(string(args[2]))
		if err != nil {
			return reply.MakeErrReply("count err")
		}
	}
	value := Set.RandomDistinctMembers(num) //随机不重复的key
	result := make([][]byte, 0)
	for i := 0; i < len(value); i++ {
		if Set.Remove(value[i]) == 1 {
			result = append(result, []byte(value[i]))
		}
	}
	return reply.MakeMultiBulkReply(result)
}

func SMEMBERS(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeArgNumErrReply("smembers")
	}
	key := args[1]
	Set, replys := db.getAsSet(string(key))
	if replys != nil {
		return replys
	}
	if Set == nil {
		return reply.MakeStatusReply("nil")
	}
	value := Set.ToSlice()
	result := make([][]byte, 0)
	for i := 0; i < len(value); i++ {
		result = append(result, []byte(value[i]))
	}
	return reply.MakeMultiBulkReply(result)
}

func SCARD(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeArgNumErrReply("scard")
	}
	key := args[1]
	Set, replys := db.getAsSet(string(key))
	if replys != nil {
		return replys
	}
	if Set == nil {
		return reply.MakeIntReply(0)
	}
	return reply.MakeIntReply(int64(Set.Len()))
}

func SREM(db *DB, args [][]byte) redis.Reply {
	if len(args) < 3 {
		return reply.MakeArgNumErrReply("srem")
	}
	key := args[1]
	Set, replys := db.getAsSet(string(key))
	if replys != nil {
		return replys
	}
	if Set == nil {
		return reply.MakeIntReply(0)
	}
	count := 0
	for i := 2; i < len(args); i++ {
		count += Set.Remove(string(args[i]))
	}
	return reply.MakeIntReply(int64(count))
}

func SRANDMEMBER(db *DB, args [][]byte) redis.Reply {
	if len(args) < 2 {
		return reply.MakeArgNumErrReply("srandmember")
	}
	key := args[1]

	Set, replys := db.getAsSet(string(key))
	if replys != nil {
		return replys
	}
	if Set == nil {
		return reply.MakeStatusReply("nil")
	}
	num := 1
	var err error
	if len(args) == 3 {
		num, err = strconv.Atoi(string(args[2]))
		if err != nil {
			return reply.MakeErrReply("count err")
		}
	}
	value := Set.ToSlice()
	result := make([][]byte, 0)
	for i := 0; i < num && i < Set.Len(); i++ {
		result = append(result, []byte(value[i]))
	}
	return reply.MakeMultiBulkReply(result)
}

func SMOVE(db *DB, args [][]byte) redis.Reply {
	if len(args) != 4 {
		return reply.MakeArgNumErrReply("smove")
	}
	key1 := args[1]
	key2 := args[2]
	num := args[3]
	Set1, replys := db.getAsSet(string(key1))
	if replys != nil {
		return replys
	}
	if Set1 == nil {
		return reply.MakeIntReply(0)
	}

	ok := Set1.Has(string(num))
	if !ok {
		return reply.MakeIntReply(0)
	}
	Set1.Remove(string(num))

	Set2, replys := db.getAsSet(string(key2))
	if replys != nil {
		return replys
	}
	if Set2 == nil {
		return reply.MakeErrReply("num2 not a set")
	}
	Set2.Add(string(num))
	return reply.MakeIntReply(1)
}

func SISMEMBER(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeArgNumErrReply("sismember")
	}
	key := args[1]
	num := args[2]
	Set, replys := db.getAsSet(string(key))
	if replys != nil {
		return replys
	}
	if Set == nil {
		return reply.MakeIntReply(0)
	}
	ok := Set.Has(string(num))
	if ok {
		return reply.MakeIntReply(1)
	}
	return reply.MakeIntReply(0)
}

func SUNION(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeArgNumErrReply("union")
	}
	key1 := args[1]
	key2 := args[2]
	result := [][]byte{}
	Set1, replys := db.getAsSet(string(key1))
	if replys != nil {
		return replys
	}
	Set2, replys := db.getAsSet(string(key2))
	if replys != nil {
		return replys
	}
	if Set1 == nil {
		result = getSetMenber(Set2)
		if result == nil {
			return reply.MakeStatusReply("(empty array)")
		}
		return reply.MakeMultiBulkReply(result)
	}
	if Set2 == nil {
		result = getSetMenber(Set1)
		if result == nil {
			return reply.MakeStatusReply("(empty array)")
		}
		return reply.MakeMultiBulkReply(result)
	}
	newset := Set1.Union(Set2)
	result = getSetMenber(newset)
	return reply.MakeMultiBulkReply(result)
}

func SUNIONSTORE(db *DB, args [][]byte) redis.Reply {
	if len(args) != 4 {
		return reply.MakeArgNumErrReply("sunionstore")
	}
	key1 := args[2]
	key2 := args[3]
	key3 := args[1]
	Set1, replys := db.getAsSet(string(key1))
	if replys != nil {
		return replys
	}
	Set2, replys := db.getAsSet(string(key2))
	if replys != nil {
		return replys
	}
	Set3, replys := db.getOrInfoSet(string(key3))
	if replys != nil {
		return replys
	}
	if Set1 == nil && Set2 == nil {
		return reply.MakeStatusReply("(empty array)")
	} else if Set1 == nil {
		Set3 = Set2
		return reply.MakeIntReply(int64(Set3.Len()))
	} else if Set2 == nil {
		Set3 = Set1
		return reply.MakeIntReply(int64(Set3.Len()))
	}
	Set3 = Set1.Union(Set2)
	data := &DataEntity{Data: Set3}
	db.Data.Put(string(key3), data)
	return reply.MakeIntReply(int64(Set3.Len()))
}

func SINTER(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeArgNumErrReply("sinter")
	}
	key1 := args[1]
	key2 := args[2]
	Set1, replys := db.getAsSet(string(key1))
	if replys != nil {
		return reply.MakeStatusReply("(empty array)")
	}
	if Set1 == nil {
		return reply.MakeStatusReply("(empty array)")
	}
	Set2, replys := db.getAsSet(string(key2))
	if replys != nil {
		return reply.MakeStatusReply("(empty array)")
	}
	if Set2 == nil {
		return reply.MakeStatusReply("(empty array)")
	}

	newset := Set1.Intersect(Set2)
	result := getSetMenber(newset)
	return reply.MakeMultiBulkReply(result)
}

func SINTERSTORE(db *DB, args [][]byte) redis.Reply {
	if len(args) != 4 {
		return reply.MakeArgNumErrReply("sinterstore")
	}
	key1 := args[2]
	key2 := args[3]
	key3 := args[1]
	Set1, replys := db.getAsSet(string(key1))
	if replys != nil {
		return reply.MakeStatusReply("(empty array)")
	}
	if Set1 == nil {
		return reply.MakeStatusReply("(empty array)")
	}
	Set2, replys := db.getAsSet(string(key2))
	if replys != nil {
		return reply.MakeStatusReply("(empty array)")
	}
	if Set2 == nil {
		return reply.MakeStatusReply("(empty array)")
	}
	Set3, replys := db.getOrInfoSet(string(key3))
	Set3 = Set1.Intersect(Set2)
	data := &DataEntity{Data: Set3}
	db.Data.Put(string(key3), data)
	return reply.MakeIntReply(int64(Set3.Len()))
}

func SDIFF(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeArgNumErrReply("sdiff")
	}
	key1 := args[1]
	key2 := args[2]
	result := [][]byte{}
	Set1, replys := db.getAsSet(string(key1))
	if replys != nil {
		return reply.MakeStatusReply("(empty array)")
	}
	Set2, replys := db.getAsSet(string(key2))
	if replys != nil {
		return reply.MakeStatusReply("(empty array)")
	}
	if Set1 == nil {
		return reply.MakeStatusReply("(empty array)")
	}
	if Set2 == nil {
		result = getSetMenber(Set1)
		if result == nil {
			return reply.MakeStatusReply("(empty array)")
		}
		return reply.MakeMultiBulkReply(result)
	}
	newset := Set1.Diff(Set2)
	result = getSetMenber(newset)
	return reply.MakeMultiBulkReply(result)
}

func SDIFFSTORE(db *DB, args [][]byte) redis.Reply {
	if len(args) != 4 {
		return reply.MakeArgNumErrReply("sdiffstore")
	}
	key1 := args[2]
	key2 := args[3]
	key3 := args[1]
	Set3, replys := db.getOrInfoSet(string(key3))
	Set1, replys := db.getAsSet(string(key1))
	if replys != nil {
		return reply.MakeStatusReply("(empty array)")
	}
	if Set1 == nil {
		return reply.MakeStatusReply("(empty array)")
	}
	Set2, replys := db.getAsSet(string(key2))
	if replys != nil {
		return reply.MakeStatusReply("(empty array)")
	}
	if Set2 == nil {
		tmp := *Set1
		Set3 = &tmp
		data := &DataEntity{Data: Set3}
		db.Data.Put(string(key3), data)
		return reply.MakeIntReply(int64(Set3.Len()))
	}
	Set3 = Set1.Diff(Set2)
	data := &DataEntity{Data: Set3}
	db.Data.Put(string(key3), data)
	return reply.MakeIntReply(int64(Set3.Len()))
}

func (db *DB) getAsSet(key string) (*set.Set, redis.Reply) {
	value, ok := db.Get(key)
	if ok == false {
		return nil, nil
	}
	bytes, ok := value.Data.(*set.Set)
	if !ok {
		return nil, &reply.WrongTypeErrReply{}
	}
	return bytes, nil
}

func (db *DB) getOrInfoSet(key string) (*set.Set, redis.Reply) {
	value, replys := db.getAsSet(key)
	if replys != nil {
		return nil, replys
	}
	if value == nil {
		newSet := set.NewSet()
		Data := &DataEntity{Data: newSet}
		db.Put(key, Data)
		//db.getAsSet(key)
		return newSet, nil
	}
	return value, nil
}

func getSetMenber(set *set.Set) [][]byte {
	if set == nil {
		return nil
	}
	value := set.ToSlice()
	result := make([][]byte, 0)
	for i := 0; i < len(value); i++ {
		result = append(result, []byte(value[i]))
	}
	return result
}
