package db

import (
	"sync"
	"test/datastruct/dict"
	"test/datastruct/list"
	"test/datastruct/lock"
	"test/interface/redis"
	"test/redis/reply"
	"time"
)

type DataEntity struct {
	Data interface{}
}

const (
	dataDictSize = 1 << 16
	ttlDictSize  = 1 << 10
	lockerSize   = 128
	aofQueueSize = 1 << 16
)

type CmdFunc func(db *DB, args [][]byte) (reply redis.Reply)

type DB struct {
	Data      dict.Dict
	TTLMap    dict.Dict
	SubMap    dict.Dict
	Locker    *lock.Locks
	interval  time.Duration
	stopWorld sync.WaitGroup //stw
}

func NewDB() *DB {
	db := &DB{
		Data:     dict.NewConcurrentDict(dataDictSize),
		TTLMap:   dict.NewConcurrentDict(ttlDictSize),
		Locker:   lock.NewLocks(lockerSize),
		interval: time.Second * 5,
	}
	/*
		AOF判断区域
	*/
	db.TimerTask()
	return db
}
func (db *DB) Exec(conn redis.Connection, args [][]byte) redis.Reply {
	//fmt.Println("我收到了相关内容")
	// for _, val := range args {
	// 	fmt.Println(string(val))
	// }
	if string(args[0]) == "COMMAND" {
		return nil
	}
	cmd := args[0]
	router := NewRouter()
	dbFunc, ok := router[string(cmd)]
	if !ok {
		return &reply.StatusReply{Status: "I donot have this command"}
	}
	reply := dbFunc(db, args)
	return reply

}

//数据处理
//获取数据
func (db *DB) Get(key string) (*DataEntity, bool) {
	db.stopWorld.Wait()
	val, ok := db.Data.Get(key)
	if ok == false {
		return nil, false
	}
	if db.IsExpired(key) {
		return nil, false
	}
	curval := val.(*DataEntity)
	return curval, true
}

//添加
func (db *DB) Put(key string, entity *DataEntity) int {
	db.stopWorld.Wait()
	return db.Data.Put(key, entity)
}

//添加如果存在
func (db *DB) PutIfExists(key string, entity *DataEntity) int {
	db.stopWorld.Wait()
	return db.Data.PutIfExists(key, entity)
}

//添加如果不存在
func (db *DB) PutIfAbsent(key string, entity *DataEntity) int {
	db.stopWorld.Wait()
	return db.Data.PutIfAbsent(key, entity)
}

//移除
func (db *DB) Remove(key string) {
	db.stopWorld.Wait()
	db.Data.Remove(key)
	db.TTLMap.Remove(key)
}

//移除多个
func (db *DB) Removes(keys ...string) (deleted int) {
	db.stopWorld.Wait()
	deleted = 0
	for _, val := range keys {
		_, exits := db.Data.Get(val)
		if exits {
			db.Data.Remove(val)
			db.TTLMap.Remove(val)
			deleted++
		}
	}
	return
}

//冲洗数据库
func (db *DB) Flush() {
	db.stopWorld.Add(1)
	defer db.stopWorld.Done()
	db.Data = dict.NewConcurrentDict(dataDictSize)
	db.TTLMap = dict.NewConcurrentDict(ttlDictSize)
	db.Locker = lock.NewLocks(lockerSize)
}

/* ---- Lock Function ----- */

func (db *DB) Lock(key string) {
	db.Locker.Lock(key)
}

func (db *DB) RLock(key string) {
	db.Locker.RLock(key)
}

func (db *DB) UnLock(key string) {
	db.Locker.UnLock(key)
}

func (db *DB) RUnLock(key string) {
	db.Locker.RUnLock(key)
}

func (db *DB) Locks(keys ...string) {
	db.Locker.Locks(keys...)
}

func (db *DB) RLocks(keys ...string) {
	db.Locker.RLocks(keys...)
}

func (db *DB) UnLocks(keys ...string) {
	db.Locker.ULocks(keys...)
}

func (db *DB) RUnLocks(keys ...string) {
	db.Locker.RUnlocks(keys...)
}

/*           TTL       */
//添加时间
func (db *DB) Expire(key string, expireTime time.Time) int {
	db.stopWorld.Wait()
	result := db.TTLMap.Put(key, expireTime)
	return result
}

//移除某个key
func (db *DB) Persist(key string) {
	db.stopWorld.Wait()
	db.TTLMap.Remove(key)
}

//执行定时器任务
func (db *DB) TimerTask() {
	ticker := time.NewTicker(db.interval)
	go func() {
		for range ticker.C {
			//fmt.Println("定时器扫描")
			db.CleanExpired()
		}
	}()
}

//清楚过期的
func (db *DB) CleanExpired() {
	now := time.Now()
	//fmt.Println(now)
	toRemove := &list.LinkedList{}
	db.TTLMap.ForEach(func(key string, val interface{}) bool {
		expireTime, _ := val.(time.Time)
		if now.After(expireTime) {
			// expired
			db.Data.Remove(key)
			toRemove.Add(key)
		}
		return true
	})
	toRemove.ForEach(func(i int, val interface{}) bool {
		key, _ := val.(string)
		db.TTLMap.Remove(key)
		return true
	})
}

//验证是否过期,如果过期就删除
func (db *DB) IsExpired(key string) bool {
	val, ok := db.TTLMap.Get(key)
	if !ok {
		return false
	}
	deletetime := val.(time.Time)
	nowtime := time.Now()
	if nowtime.After(deletetime) {
		db.TTLMap.Remove(key)
		return true
	} else {
		return false
	}
}
