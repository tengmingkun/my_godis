package db

import (
	"strconv"
	"test/datastruct/zset"
	"test/interface/redis"
	"test/redis/reply"
)

func ZADD(db *DB, args [][]byte) redis.Reply {
	if len(args)%2 != 0 {
		return reply.MakeArgNumErrReply("zadd")
	}
	key := args[1]
	zset, replys := db.getOrInfoZset(string(key))
	if replys != nil {
		return replys
	}
	count := 0
	for i := 2; i < len(args); i += 2 {
		num, err := strconv.Atoi(string(args[i]))
		if err != nil {
			return reply.MakeErrReply("num has err")
		}
		val := args[i+1]
		zset.ADD(num, val)
		count++
	}
	return reply.MakeIntReply(int64(count))
}

func ZCARD(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeArgNumErrReply("zcard")
	}
	key := args[1]
	zset, replys := db.getAsZset(string(key))
	if replys != nil {
		return replys
	}
	if zset == nil {
		return reply.MakeIntReply(0)
	}
	return reply.MakeIntReply(int64(zset.Size))
}

func ZRANGE(db *DB, args [][]byte) redis.Reply {
	if len(args) < 4 {
		return reply.MakeArgNumErrReply("zrange")
	}
	key := args[1]
	start, err := strconv.Atoi(string(args[2]))
	if err != nil {
		return reply.MakeErrReply("string to num error ")
	}
	end, err := strconv.Atoi(string(args[3]))
	if err != nil {
		return reply.MakeErrReply("string to num error ")
	}
	zset, replys := db.getAsZset(string(key))
	if replys != nil {
		return replys
	}
	if zset == nil {
		return reply.MakeStatusReply("nil")
	}
	socre, value := zset.GetRange(start, end)
	var bytes [][]byte
	if len(args) == 5 {
		bytes = make([][]byte, 2*len(socre))
		i := 0
		for _, v := range value {
			bytes[i] = v.([]byte)
			i = i + 2
		}
		i = 1
		for _, v := range socre {
			num := strconv.Itoa(v)
			bytes[i] = []byte(num)
			i = i + 2
		}
	} else {
		for _, v := range value {
			bytes = append(bytes, v.([]byte))
		}
	}
	return reply.MakeMultiBulkReply(bytes)
}

func ZREVRANGE(db *DB, args [][]byte) redis.Reply {
	if len(args) < 4 {
		return reply.MakeArgNumErrReply("zrevrange")
	}
	key := args[1]
	start, err := strconv.Atoi(string(args[2]))
	if err != nil {
		return reply.MakeErrReply("string to num error ")
	}
	end, err := strconv.Atoi(string(args[3]))
	if err != nil {
		return reply.MakeErrReply("string to num error ")
	}
	zset, replys := db.getAsZset(string(key))
	if replys != nil {
		return replys
	}
	if zset == nil {
		return reply.MakeStatusReply("nil")
	}
	socre, value := zset.GetRange(start, end)
	var bytes [][]byte
	if len(args) == 5 {
		bytes = make([][]byte, 2*len(socre))
		i := 0
		for j := len(value) - 1; j >= 0; j-- {
			bytes[i] = value[j].([]byte)
			i = i + 2
		}
		i = 1
		for j := len(socre) - 1; j >= 0; j-- {
			num := strconv.Itoa(socre[j])
			bytes[i] = []byte(num)
			i = i + 2
		}
	} else {
		for j := len(value) - 1; j >= 0; j-- {
			bytes = append(bytes, value[j].([]byte))
		}
	}
	return reply.MakeMultiBulkReply(bytes)
}

func ZCOUNT(db *DB, args [][]byte) redis.Reply {
	if len(args) != 4 {
		return reply.MakeArgNumErrReply("zcount")
	}
	key := args[1]
	mins, err := strconv.Atoi(string(args[2]))
	if err != nil {
		return reply.MakeErrReply("string to num error")
	}
	maxs, err := strconv.Atoi(string(args[3]))
	if err != nil {
		return reply.MakeErrReply("string to num error")
	}
	zset, replys := db.getAsZset(string(key))
	if replys != nil {
		return replys
	}
	if zset == nil {
		return reply.MakeStatusReply("nil")
	}
	count := zset.GetCount(mins, maxs)
	return reply.MakeIntReply(int64(count))
}

func ZSCORE(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeArgNumErrReply("zsore")
	}
	key := args[1]
	member := args[2]
	zset, replys := db.getAsZset(string(key))
	if replys != nil {
		return replys
	}
	if zset == nil {
		return reply.MakeStatusReply("(nil)")
	}
	result, ok := zset.FindMember(member)
	if !ok {
		return reply.MakeStatusReply("(nil)")
	}
	return reply.MakeIntReply(int64(result))
}

func ZSCAN(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeArgNumErrReply("zscan")
	}
	return nil
}

func ZRANK(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeArgNumErrReply("zrank")
	}
	key := args[1]
	member := args[2]
	zset, replys := db.getAsZset(string(key))
	if replys != nil {
		return replys
	}
	if zset == nil {
		return reply.MakeIntReply(-1)
	}
	count := zset.GetRank(member)
	return reply.MakeIntReply(int64(count))
}

func (db *DB) getAsZset(key string) (*zset.Zset, redis.Reply) {
	val, ok := db.Get(key)
	if ok == false {
		return nil, nil
	}
	bytes := val.Data.(*zset.Zset)
	if bytes == nil {
		return nil, &reply.WrongTypeErrReply{}
	}
	return bytes, nil
}
func (db *DB) getOrInfoZset(key string) (*zset.Zset, redis.Reply) {
	value, replys := db.getAsZset(key)
	if replys != nil {
		return nil, replys
	}
	if value == nil {
		newzset := zset.NewZset()
		data := &DataEntity{Data: newzset}
		db.Put(key, data)
		return newzset, nil
	}
	return value, nil
}
