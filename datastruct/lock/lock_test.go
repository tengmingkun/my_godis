package lock

import (
	"fmt"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	lock := NewLocks(20)
	for i := 0; i < 10; i++ {
		go GetLock(lock, "test")
	}
	time.Sleep(time.Second * 20)
}
func GetLock(lock *Locks, key string) {
	lock.Lock(key)
	fmt.Println("我被访问了")
	fmt.Println("我拿到了锁")
	time.Sleep(time.Second * 2)
	lock.UnLock(key)
}
