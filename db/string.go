package db

import (
	"test/interface/redis"
	"test/redis/reply"
)

func Get(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeArgNumErrReply("get")
	}
	key := string(args[1])
	value, ok := db.Get(key)
	if !ok {
		return reply.MakeStatusReply("(nil)")
	}
	return reply.MakeBulkReply(value.Data.([]byte))
}
func GetSet(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeArgNumErrReply("get")
	}
	key := string(args[1])
	value1, ok := db.Get(key)
	value2 := &DataEntity{Data: args[2]}
	db.Put(string(key), value2)
	if !ok {
		return reply.MakeStatusReply("(nil)")
	}
	return reply.MakeBulkReply(value1.Data.([]byte))

}
func StrLen(db *DB, args [][]byte) redis.Reply {
	value, ok := db.Get(string(args[1]))
	if !ok {
		return reply.MakeIntReply(0)
	} else {
		length := len(string(value.Data.([]byte)))
		return reply.MakeIntReply(int64(length))
	}
}

func SET(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		reply := reply.MakeArgNumErrReply("set")
		return reply
	}
	key := args[1]
	value := &DataEntity{Data: args[2]}
	db.Put(string(key), value)
	return &reply.StatusReply{Status: "ok"}
}

func MSET(db *DB, args [][]byte) redis.Reply {
	if len(args)%2 == 0 {
		reply := reply.MakeArgNumErrReply("mset")
		return reply
	}
	for i := 1; i < len(args); i += 2 {
		key := args[i]
		value := &DataEntity{Data: args[i+1]}
		db.Put(string(key), value)
	}
	return reply.MakeStatusReply("ok")
}
func SetNx(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeArgNumErrReply("setnx")
	}
	key := string(args[1])
	value := &DataEntity{
		Data: args[2],
	}
	db.PutIfAbsent(key, value)
	return reply.MakeStatusReply("OK")
}

func SetEx(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeArgNumErrReply("setex")
	}
	key := string(args[1])
	value := &DataEntity{Data: args[2]}
	db.PutIfExists(key, value)
	return reply.MakeStatusReply("OK")
}
