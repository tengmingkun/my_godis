package db

func NewRouter() map[string]CmdFunc {
	router := make(map[string]CmdFunc)
	//string 相关
	router["set"] = SET
	router["mset"] = MSET
	router["setnx"] = SetNx
	router["setex"] = SetEx
	router["get"] = Get
	router["strlen"] = StrLen
	router["getset"] = GetSet

	//key相关
	router["expire"] = Expire
	router["pexpireat"] = PExpireAt
	router["ttl"] = TTL
	router["rename"] = ReName
	router["del"] = DEL

	//连接相关
	router["echo"] = ECHO
	router["ping"] = PING

	//list相关
	router["rpush"] = RPUSH
	router["lindex"] = LINDEX
	router["lrange"] = LRANGE
	router["lpoprpush"] = LPOPRPUSH
	router["lrem"] = LREM
	router["llen"] = LLEN
	router["lpop"] = LPOP
	router["lpushx"] = LPUSHX
	router["lpush"] = LPUSH
	router["lset"] = LSET
	router["rpop"] = RPOP
	router["rpushx"] = RPUSHX

	//set 相关
	router["sadd"] = SADD
	router["smembers"] = SMEMBERS
	router["scard"] = SCARD
	router["srem"] = SREM
	router["spop"] = SPOP
	router["srandmember"] = SRANDMEMBER
	router["smove"] = SMOVE
	router["sismember"] = SISMEMBER
	router["sunion"] = SUNION
	router["sinter"] = SINTER
	router["sdiff"] = SDIFF
	router["sunionstore"] = SUNIONSTORE
	router["sinterstore"] = SINTERSTORE
	router["sdiffstore"] = SDIFFSTORE
	return router

}
