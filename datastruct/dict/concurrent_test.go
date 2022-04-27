package dict

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func TestNewDict(t *testing.T) {
	dict := NewConcurrentDict(23)
	if dict == nil {
		t.Error("New dict error")
	}
}

func TestPutAndGetAndLenAndRemove(t *testing.T) {
	dict := NewConcurrentDict(23)
	dict.Put("name", "tengmk")
	val, statu := dict.Get("name")
	if val != "tengmk" || statu != true {
		t.Error("get and put error")
	}
	len := dict.Len()
	if len != 1 {
		t.Error("len error")
	}
	dict.Remove("name")
	dict.Put("age", 25)
	val, statu = dict.Get("name")
	if val != nil || statu != false {
		t.Error("remove error")
	}
	len = dict.Len()
	fmt.Println(len)
	if len != 1 {
		t.Error("len error")
	}
}

func TestPutIfAbset(t *testing.T) {
	dict := NewConcurrentDict(23)
	statu := dict.PutIfAbsent("name", "tengmk")
	if statu != 1 {
		t.Error("putifabsen error")
	}
	dict.Put("name", "tengmk")
	statu = dict.PutIfAbsent("name", "zhangxy")
	if statu != 0 {
		t.Error("put if absen error")
	}
	val, _ := dict.Get("name")
	fmt.Println(val)
	if val != "tengmk" {
		t.Error("pus if absen error")
	}
}
func TestPutIfexits(t *testing.T) {
	dict := NewConcurrentDict(23)
	statu := dict.PutIfExists("name", "tengmk")
	fmt.Println("put if exits statu", statu)
	if statu != 0 {
		t.Error("put exits error")
	}
	dict.Put("name", "tengmk")
	statu = dict.PutIfExists("name", "zhangxy")
	if statu != 1 {
		t.Error("put exits error")
	}
	val, _ := dict.Get("name")
	fmt.Println(val)
	if val != "zhangxy" {
		t.Error("put exits error")
	}
}

func TestPut(t *testing.T) {
	d := NewConcurrentDict(0)
	count := 100
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			// insert
			key := "k" + strconv.Itoa(i)
			ret := d.Put(key, i)
			if ret != 1 { // insert 1
				t.Error("put test failed: expected result 1, actual: " + strconv.Itoa(ret) + ", key: " + key)
			}
			val, ok := d.Get(key)
			if ok {
				intVal, _ := val.(int)
				if intVal != i {
					t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal) + ", key: " + key)
				}
			} else {
				_, ok := d.Get(key)
				t.Error("put test failed: expected true, actual: false, key: " + key + ", retry: " + strconv.FormatBool(ok))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
func TestPutIfAbsent(t *testing.T) {
	d := NewConcurrentDict(0)
	count := 100
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			// insert
			key := "k" + strconv.Itoa(i)
			ret := d.PutIfAbsent(key, i)
			if ret != 1 { // insert 1
				t.Error("put test failed: expected result 1, actual: " + strconv.Itoa(ret) + ", key: " + key)
			}
			val, ok := d.Get(key)
			if ok {
				intVal, _ := val.(int)
				if intVal != i {
					t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal) +
						", key: " + key)
				}
			} else {
				_, ok := d.Get(key)
				t.Error("put test failed: expected true, actual: false, key: " + key + ", retry: " + strconv.FormatBool(ok))
			}

			// update
			ret = d.PutIfAbsent(key, i*10)
			if ret != 0 { // no update
				t.Error("put test failed: expected result 0, actual: " + strconv.Itoa(ret))
			}
			val, ok = d.Get(key)
			if ok {
				intVal, _ := val.(int)
				if intVal != i {
					t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal) + ", key: " + key)
				}
			} else {
				t.Error("put test failed: expected true, actual: false, key: " + key)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestPutIfExists(t *testing.T) {
	d := NewConcurrentDict(0)
	count := 100
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			// insert
			key := "k" + strconv.Itoa(i)
			// insert
			ret := d.PutIfExists(key, i)
			if ret != 0 { // insert
				t.Error("put test failed: expected result 0, actual: " + strconv.Itoa(ret))
			}

			d.Put(key, i)
			ret = d.PutIfExists(key, 10*i)
			val, ok := d.Get(key)
			if ok {
				intVal, _ := val.(int)
				if intVal != 10*i {
					t.Error("put test failed: expected " + strconv.Itoa(10*i) + ", actual: " + strconv.Itoa(intVal))
				}
			} else {
				_, ok := d.Get(key)
				t.Error("put test failed: expected true, actual: false, key: " + key + ", retry: " + strconv.FormatBool(ok))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestRemove(t *testing.T) {
	d := NewConcurrentDict(0)

	// remove head node
	for i := 0; i < 100; i++ {
		// insert
		key := "k" + strconv.Itoa(i)
		d.Put(key, i)
	}
	for i := 0; i < 100; i++ {
		key := "k" + strconv.Itoa(i)

		val, ok := d.Get(key)
		if ok {
			intVal, _ := val.(int)
			if intVal != i {
				t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal))
			}
		} else {
			t.Error("put test failed: expected true, actual: false")
		}

		ret := d.Remove(key)
		if ret != 1 {
			t.Error("remove test failed: expected result 1, actual: " + strconv.Itoa(ret) + ", key:" + key)
		}
		_, ok = d.Get(key)
		if ok {
			t.Error("remove test failed: expected true, actual false")
		}
		ret = d.Remove(key)
		if ret != 0 {
			t.Error("remove test failed: expected result 0 actual: " + strconv.Itoa(ret))
		}
	}

	// remove tail node
	d = NewConcurrentDict(0)
	for i := 0; i < 100; i++ {
		// insert
		key := "k" + strconv.Itoa(i)
		d.Put(key, i)
	}
	for i := 9; i >= 0; i-- {
		key := "k" + strconv.Itoa(i)

		val, ok := d.Get(key)
		if ok {
			intVal, _ := val.(int)
			if intVal != i {
				t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal))
			}
		} else {
			t.Error("put test failed: expected true, actual: false")
		}

		ret := d.Remove(key)
		if ret != 1 {
			t.Error("remove test failed: expected result 1, actual: " + strconv.Itoa(ret))
		}
		_, ok = d.Get(key)
		if ok {
			t.Error("remove test failed: expected true, actual false")
		}
		ret = d.Remove(key)
		if ret != 0 {
			t.Error("remove test failed: expected result 0 actual: " + strconv.Itoa(ret))
		}
	}

	// remove middle node
	d = NewConcurrentDict(0)
	d.Put("head", 0)
	for i := 0; i < 10; i++ {
		// insert
		key := "k" + strconv.Itoa(i)
		d.Put(key, i)
	}
	d.Put("tail", 0)
	for i := 9; i >= 0; i-- {
		key := "k" + strconv.Itoa(i)

		val, ok := d.Get(key)
		if ok {
			intVal, _ := val.(int)
			if intVal != i {
				t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal))
			}
		} else {
			t.Error("put test failed: expected true, actual: false")
		}

		ret := d.Remove(key)
		if ret != 1 {
			t.Error("remove test failed: expected result 1, actual: " + strconv.Itoa(ret))
		}
		_, ok = d.Get(key)
		if ok {
			t.Error("remove test failed: expected true, actual false")
		}
		ret = d.Remove(key)
		if ret != 0 {
			t.Error("remove test failed: expected result 0 actual: " + strconv.Itoa(ret))
		}
	}
}

func TestForEach(t *testing.T) {
	d := NewConcurrentDict(0)
	size := 100
	for i := 0; i < size; i++ {
		// insert
		key := "k" + strconv.Itoa(i)
		d.Put(key, i)
	}
	i := 0
	d.ForEach(func(key string, value interface{}) bool {
		intVal, _ := value.(int)
		expectedKey := "k" + strconv.Itoa(intVal)
		if key != expectedKey {
			t.Error("remove test failed: expected " + expectedKey + ", actual: " + key)
		}
		i++
		return true
	})
	if i != size {
		t.Error("remove test failed: expected " + strconv.Itoa(size) + ", actual: " + strconv.Itoa(i))
	}
}
