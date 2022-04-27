package db

import (
	"test/interface/redis"
	"test/redis/reply"
)

func ECHO(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeArgNumErrReply("ECHO")
	}
	str := args[1]
	return reply.MakeBulkReply(str)
}

func PING(db *DB, args [][]byte) redis.Reply {
	return reply.MakeStatusReply("PONG")
}
