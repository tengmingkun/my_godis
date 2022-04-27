package db

import (
	"fmt"
	"strconv"
	"test/interface/redis"
	"test/redis/reply"
	"time"
)

func Expire(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'pexpireat' command")
	}
	key := string(args[1])

	raw, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("ERR value is not an integer or out of range")
	}
	fmt.Println(raw)
	expireTime := time.Now().Add(time.Second * time.Duration(raw))
	fmt.Println(expireTime)
	_, exists := db.Get(key)
	if !exists {
		return reply.MakeIntReply(0)
	}

	result := db.Expire(key, expireTime)
	return reply.MakeIntReply(int64(result))
}

func PExpireAt(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'pexpireat' command")
	}
	key := string(args[1])

	raw, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("ERR value is not an integer or out of range")
	}
	expireTime := time.Unix(0, raw*int64(time.Millisecond))
	_, exists := db.Get(key)
	if !exists {
		return reply.MakeIntReply(0)
	}

	result := db.Expire(key, expireTime)
	return reply.MakeIntReply(int64(result))
}

func TTL(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeArgNumErrReply("ttl")
	}
	key := string(args[1])
	_, ok := db.Get(key)
	if ok == false {
		return reply.MakeIntReply(-2)
	}
	value, ok := db.TTLMap.Get(key)
	if ok == false {
		return reply.MakeIntReply(-1)
	} else {
		lastTime := value.(time.Time).Sub(time.Now())
		return reply.MakeIntReply(int64(lastTime) / 1000000000)
	}
}

func EXISTS(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeArgNumErrReply("exists")
	}
	key := string(args[1])
	_, ok := db.Get(key)
	if !ok {
		return reply.MakeIntReply(0)
	} else {
		return reply.MakeIntReply(1)
	}
}

func ReName(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeArgNumErrReply("rename")
	}
	oldkey := string(args[1])
	newKey := string(args[2])
	if oldkey == newKey {
		return reply.MakeErrReply("old key is the same as old key")
	}
	value, ok := db.Get(oldkey)
	if !ok {
		return reply.MakeErrReply("no such key")
	} else {
		key_time, ok := db.TTLMap.Get(oldkey)
		db.Remove(oldkey)
		db.Put(newKey, value)
		if ok {
			db.TTLMap.Put(newKey, key_time)
		}
	}
	return reply.MakeStatusReply("ok")
}

func DEL(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeArgNumErrReply("del")
	}
	key := string(args[1])
	count := db.Removes(key)
	return reply.MakeIntReply(int64(count))
}
