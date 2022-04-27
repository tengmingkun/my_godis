package list

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	list := NewLinkedList()
	// list.Add(123)
	// if list.first.val != 123 || list.last.val != 123 {
	// 	fmt.Println("添加失败")
	// }

	for i := 0; i < 10; i++ {
		list.Add(i)
	}
	i := 0
	n := list.first
	for n != nil {
		if n.val != i {
			t.Error("add error")
		}
		n = n.next
		i++
	}
}

func TestFind(t *testing.T) {
	list := NewLinkedList()
	for i := 0; i < 10; i++ {
		list.Add(i)
	}
	n := list.find(3)
	if n.val != 3 {
		t.Error("find error")
	}
}
func TestGet(t *testing.T) {
	list := NewLinkedList()
	for i := 0; i < 10; i++ {
		list.Add(i)
	}
	n := list.Get(4)
	if n != 4 {
		t.Error("Get error")
	}
}

func TestSet(t *testing.T) {
	list := NewLinkedList()
	for i := 0; i < 10; i++ {
		list.Add(i)
	}
	list.Set(4, 44)
	val := list.Get(4)
	if val != 44 {
		t.Error("set error")
	}
}
func TestInsert(t *testing.T) {
	list := NewLinkedList()
	for i := 0; i < 10; i++ {
		list.Add(i)
	}
	list.Insert(3, 33)
	val := list.Get(3)
	if val != 33 {
		fmt.Println(val)
		t.Error("insert error")
	}
	list.Insert(0, 33)
	val = list.Get(0)
	if val != 33 {
		t.Error("insert error")
	}
	list.Insert(list.size-1, 33)
	val = list.Get(list.size - 1)
	if val != 33 {
		t.Error("insert error")
	}

}

func TestRemove(t *testing.T) {
	list := NewLinkedList()
	for i := 0; i < 10; i++ {
		list.Add(i)
	}
	n := list.Remove(5)
	if n != 5 {
		t.Error("remove error")
	}
	list.Remove(0)
	if list.first.val != 1 {
		t.Error("remove error")
	}
	list.Remove(list.size - 1)
	if list.last.val != 8 {
		t.Error("remove error")
	}
}

func TestRemoveLast(t *testing.T) {
	list := NewLinkedList()
	for i := 0; i < 10; i++ {
		list.Add(i)
	}
	list.RemoveLast()
	if list.last.val != 8 {
		t.Error("remove last error")
	}
	list1 := NewLinkedList()
	list1.Add(23)
	list1.RemoveLast()
	if list1.first != nil || list1.last != nil {
		t.Error("removelast error")
	}
}

func TestRemoveAllByVal(t *testing.T) {
	list := NewLinkedList()
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			list.Add(10)
		} else {
			list.Add(5)
		}

	}
	n := list.RemoveAllByVal(10)
	if list.size != 5 || n != 5 {
		t.Error("removeallbyval error")
	}

}
