package db

import (
	"strconv"
	"test/datastruct/list"
	"test/interface/redis"
	"test/redis/reply"
)

func LINDEX(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeArgNumErrReply("Lindex")
	}
	key := string(args[1])
	index, err := strconv.Atoi(string(args[2]))
	if err != nil {
		return reply.MakeErrReply("Err index err")
	}
	List, replys := db.getAsList(key)
	if replys != nil {
		return replys
	}
	if index >= List.Len() {
		return reply.MakeErrReply("index out of range ")
	}
	value := List.Get(index)
	return reply.MakeBulkReply([]byte(value.([]byte)))
}

func RPUSH(db *DB, args [][]byte) redis.Reply {
	if len(args) == 1 {
		return reply.MakeArgNumErrReply("rpush")
	}
	key := string(args[1])
	//fmt.Println("len of rpush ", len(args))
	oldList, replys := db.getOrInfoList(key)
	if replys != nil {
		return replys
	}
	for i := 2; i < len(args); i++ {
		oldList.Add(args[i])
	}
	return reply.MakeIntReply(int64(oldList.Len()))
}

func LRANGE(db *DB, args [][]byte) redis.Reply {
	if len(args) != 4 {
		return reply.MakeArgNumErrReply("Lrange")
	}
	key := string(args[1])
	pos1, err := strconv.Atoi(string(args[2]))
	if err != nil {
		return reply.MakeErrReply("Err index err")
	}
	pos2, err := strconv.Atoi(string(args[3]))
	if err != nil {
		return reply.MakeErrReply("Err index err")
	}
	List, replys := db.getAsList(key)
	if replys != nil {
		return replys
	}
	if List == nil {
		return reply.MakeMultiBulkReply(nil)
	}
	//fmt.Println(pos1, "   ", pos2)
	if pos2 < pos1 {
		return reply.MakeMultiBulkReply(nil)
	}
	if pos1 < 0 {
		pos1 = 0
	}
	if pos2 >= List.Len() {
		pos2 = List.Len() - 1
	}
	var result [][]byte
	for i := pos1; i <= pos2; i++ {
		value := List.Get(i)
		result = append(result, value.([]byte))
	}
	return reply.MakeMultiBulkReply(result)
}

func LPOPRPUSH(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeArgNumErrReply("lpoprpush")
	}
	key1 := args[1]
	key2 := args[2]
	value1, ok := db.Get(string(key1))
	if !ok {
		return reply.MakeStatusReply("nil")
	}
	list1, ok := value1.Data.(*list.LinkedList)
	if !ok {
		return reply.MakeErrReply(string(key1) + "is not a list")
	}
	if list1.Len() == 0 {
		return reply.MakeStatusReply("nil")
	}
	num := list1.Get(list1.Len() - 1)
	list1.Remove(list1.Len() - 1)
	list2, replys := db.getOrInfoList(string(key2))
	if replys != nil {
		return replys
	}
	list2.Add(num)
	return reply.MakeBulkReply(num.([]byte))
}
func LREM(db *DB, args [][]byte) redis.Reply {
	if len(args) != 4 {
		return reply.MakeArgNumErrReply("lrem")
	}
	key := string(args[1])
	List, _ := db.getAsList(key)
	if List == nil {
		return reply.MakeStatusReply("nil")
	}
	count, err := strconv.Atoi(string(args[2]))
	if err != nil {
		return reply.MakeErrReply("count err")
	}
	target := args[3]
	var num int
	if count > 0 {
		num = List.RemoveByval(target, count)
	} else if count < 0 {
		num = List.ReverseRemoveByVal(target, -count)
	} else {
		num = List.RemoveAllByVal(target)
	}
	return reply.MakeIntReply(int64(num))
}

func LLEN(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeArgNumErrReply("llen")
	}
	key := args[1]
	List, replys := db.getAsList(string(key))
	if replys != nil {
		return replys
	}
	if List == nil {
		return reply.MakeStatusReply("nil")
	}
	len := List.Len()
	return reply.MakeIntReply(int64(len))

}

func LPOP(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeArgNumErrReply("lpop")
	}
	key := args[1]
	List, replys := db.getAsList(string(key))
	if replys != nil {
		return replys
	}
	if List == nil || List.Len() == 0 {
		return reply.MakeStatusReply("nil")
	}
	value := List.Remove(0)
	return reply.MakeBulkReply(value.([]byte))
}

func LPUSHX(db *DB, args [][]byte) redis.Reply {
	if len(args) < 2 {
		return reply.MakeArgNumErrReply("lpushx")
	}
	key := args[1]
	List, replys := db.getAsList(string(key))
	if replys != nil {
		return replys
	}
	if List == nil {
		return reply.MakeIntReply(0)
	}
	for i := 2; i < len(args); i++ {
		value := args[i]
		List.Insert(0, value)
	}
	return reply.MakeIntReply(int64(List.Len()))
}

func LPUSH(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return reply.MakeArgNumErrReply("push")
	}
	key := args[1]
	List, replys := db.getOrInfoList(string(key))
	if replys != nil {
		return replys
	}
	value := args[2]
	//fmt.Println(value)
	List.Insert(0, value)
	return reply.MakeIntReply(int64(List.Len()))
}

func LSET(db *DB, args [][]byte) redis.Reply {
	if len(args) != 4 {
		return reply.MakeArgNumErrReply("lset")
	}
	key := args[1]
	index, err := strconv.Atoi(string(args[2]))
	value := args[3]
	if err != nil {
		return reply.MakeErrReply("index is err")
	}
	List, replys := db.getAsList(string(key))
	if replys != nil {
		return replys
	}
	if List == nil {
		return reply.MakeErrReply("List is nil")
	}
	if index < 0 || index >= List.Len() {
		return reply.MakeErrReply("index out of range ")
	}
	List.Set(index, value)
	return reply.MakeStatusReply("OK")
}

func RPOP(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeArgNumErrReply("rpop")
	}
	key := args[1]
	List, replys := db.getAsList(string(key))
	if replys != nil {
		return replys
	}
	if List == nil {
		return reply.MakeStatusReply("nil")
	}
	if List.Len() == 0 {
		return reply.MakeStatusReply("nil")
	}
	value := List.Remove(List.Len() - 1)
	return reply.MakeBulkReply([]byte(value.([]byte)))
}

func RPUSHX(db *DB, args [][]byte) redis.Reply {
	if len(args) <= 2 {
		return reply.MakeArgNumErrReply("rpushx")
	}
	key := args[1]
	List, replys := db.getAsList(string(key))
	if replys != nil {
		return replys
	}
	if List == nil {
		return reply.MakeIntReply(0)
	}
	for i := 2; i < len(args); i++ {
		value := args[i]
		List.Insert(List.Len(), value)
	}
	return reply.MakeIntReply(int64(List.Len()))
}

func (db *DB) getAsList(key string) (*list.LinkedList, redis.Reply) {
	value, ok := db.Get(key)
	if ok == false {
		return nil, nil
	}
	bytes, ok := value.Data.(*list.LinkedList)

	if !ok {
		return nil, &reply.WrongTypeErrReply{}
	}
	return bytes, nil
}

func (db *DB) getOrInfoList(key string) (*list.LinkedList, redis.Reply) {
	value, reply := db.getAsList(key)
	if reply != nil {
		return nil, reply
	}
	if value == nil {
		newList := list.MakeBytesList()
		data := &DataEntity{Data: newList}
		db.Put(key, data)
		return newList, nil
	}
	return value, nil
}
