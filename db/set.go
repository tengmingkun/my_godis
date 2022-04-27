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
