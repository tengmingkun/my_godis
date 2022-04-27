package list

type LinkedList struct {
	first *node
	last  *node
	size  int
}
type node struct {
	val  interface{}
	prev *node
	next *node
}

func NewLinkedList() *LinkedList {
	return &LinkedList{size: 0}
}
func (list *LinkedList) Add(val interface{}) {
	if list == nil {
		panic("list is nil")
	}
	p := &node{val: val, prev: nil, next: nil}
	if list.last == nil {
		list.first = p
		list.last = p
	} else {
		p.prev = list.last
		list.last.next = p
		list.last = p
	}
	list.size++
}

func (list *LinkedList) find(index int) (n *node) {
	if list == nil {
		panic("list is nil")
	}
	n = list.first
	for index > 0 {
		n = n.next
		index--
	}
	return n
}

func (list *LinkedList) Get(index int) (val interface{}) {
	if list == nil {
		panic("list is nil")
	}
	if index < 0 || index >= list.size {
		panic("index out of bound")
	}
	return list.find(index).val
}

func (list *LinkedList) Set(index int, val interface{}) {
	if list == nil {
		panic("list is nil")
	}
	if index < 0 || index >= list.size {
		panic("index out of bound")
	}
	n := list.find(index)
	n.val = val
}

func (list *LinkedList) Insert(index int, val interface{}) {
	if list == nil {
		panic("list is nil")
	}
	if index < 0 || index > list.size {
		panic("index out of bound")
	}

	if index == list.size {
		list.Add(val)
		return
	} else {
		// list is not empty
		pivot := list.find(index)
		n := &node{
			val:  val,
			prev: pivot.prev,
			next: pivot,
		}
		if pivot.prev == nil {
			list.first = n
		} else {
			pivot.prev.next = n
		}
		pivot.prev = n
		list.size++
	}
}
func (list *LinkedList) removeNode(n *node) {
	if n.prev == nil {
		list.first = n.next
	} else {
		n.prev.next = n.next
	}
	if n.next == nil {
		list.last = n.prev
	} else {
		n.next.prev = n.prev
	}
	n.prev = nil
	n.next = nil
	list.size--
}

func (list *LinkedList) Remove(index int) (val interface{}) {
	if index < 0 || index >= list.size {
		return
	}
	n := list.find(index)
	val = n.val
	list.removeNode(n)
	return
}

func (list *LinkedList) RemoveLast() (val interface{}) {
	if list.last == nil {
		return nil
	}
	n := list.last
	val = n.val
	list.removeNode(n)
	return
}

func Equals(a interface{}, b interface{}) bool {
	sliceA, okA := a.([]byte)
	sliceB, okB := b.([]byte)
	if okA && okB {
		return BytesEquals(sliceA, sliceB)
	}
	return a == b
}

func BytesEquals(a []byte, b []byte) bool {
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	size := len(a)
	for i := 0; i < size; i++ {
		av := a[i]
		bv := b[i]
		if av != bv {
			return false
		}
	}
	return true
}
func (list *LinkedList) RemoveAllByVal(val interface{}) int {
	n := list.first
	removed := 0
	for n != nil {
		var toRemoveNode *node
		if Equals(n.val, val) {
			toRemoveNode = n
		}
		if n.next == nil {
			if toRemoveNode != nil {
				removed++
				list.removeNode(toRemoveNode)
			}
			break
		} else {
			n = n.next
		}
		if toRemoveNode != nil {
			removed++
			list.removeNode(toRemoveNode)
		}
	}
	return removed
}

func (list *LinkedList) Len() int {
	return list.size
}

func (list *LinkedList) RemoveByval(val interface{}, count int) int {
	n := list.first
	removed := 0
	for n != nil {
		var toRemoveNode *node
		if Equals(n.val, val) {
			toRemoveNode = n
		}
		if n.next == nil {
			if toRemoveNode != nil {
				removed++
				list.removeNode(toRemoveNode)
			}
			break
		} else {
			n = n.next
		}
		if toRemoveNode != nil {
			removed++
			list.removeNode(toRemoveNode)

		}
		if removed == count {
			break
		}
	}
	return removed
}

func (list *LinkedList) ReverseRemoveByVal(val interface{}, count int) int {
	if list == nil {
		panic("list is nil")
	}
	n := list.last
	removed := 0
	for n != nil {
		var toRemoveNode *node
		if Equals(n.val, val) {
			toRemoveNode = n
		}
		if n.prev == nil {
			if toRemoveNode != nil {
				removed++
				list.removeNode(toRemoveNode)
			}
			break
		} else {
			n = n.prev
		}

		if toRemoveNode != nil {
			removed++
			list.removeNode(toRemoveNode)
		}
		if removed == count {
			break
		}
	}
	return removed
}

func (list *LinkedList) ForEach(consumer func(int, interface{}) bool) {
	if list == nil {
		return
	}
	n := list.first
	i := 0
	for n != nil {
		gonext := consumer(i, n.val)
		if !gonext || n.next == nil {
			break
		} else {
			i++
			n = n.next
		}

	}
}
func (list *LinkedList) Contains(val interface{}) bool {
	contains := false
	list.ForEach(func(i int, actual interface{}) bool {
		if actual == val {
			contains = true
			return false
		}
		return true
	})
	return contains
}

func MakeBytesList(vals ...[]byte) *LinkedList {
	list := LinkedList{}
	for _, v := range vals {
		list.Add(v)
	}
	return &list
}
