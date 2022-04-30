package skiplist

import (
	"math/rand"
	"test/datastruct/stack"
)

const MAX_LEVEL = 32
const INT_MAX = int(^uint(0) >> 1)
const INT_MIN = ^INT_MAX

type Node struct {
	key   int
	val   interface{}
	right *Node
	down  *Node
}

type SkipList struct {
	HeadNode  *Node
	HighLevel int
}

func NewNode(key int, val interface{}) *Node {
	return &Node{key: key, val: val}
}

func NewSkipList() *SkipList {
	return &SkipList{HeadNode: NewNode(INT_MIN, 0), HighLevel: 0}
}

func (s *SkipList) Search(key int) *Node {
	node := s.HeadNode
	for node != nil {
		if node.key == key {
			return node
		} else if node.right == nil {
			node = node.down
		} else if node.right.key > key {
			node = node.down
		} else {
			node = node.right
		}
	}
	return nil
}

func (s *SkipList) Delete(key int) bool {
	node := s.HeadNode
	for node != nil {
		if node.right == nil {
			node = node.down
		} else if node.right.key == key {
			node.right = node.right.right
			node = node.down
		} else if node.right.key > key {
			node = node.down
		} else {
			node = node.right
		}
	}
	return true
}

func (s *SkipList) Update(key int, val interface{}) {
	if node := s.Search(key); node != nil {
		for node != nil {
			node.val = val
			node = node.down
		}
		return
	}

}

func (s *SkipList) Insert(key int, val interface{}) {
	node := s.Search(key)
	for node != nil {
		node.val = val
		node = node.down
	}
	stacks := stack.NewStack()
	node = s.HeadNode
	for node != nil {
		if node.right == nil {
			stacks.Push(node)
			node = node.down
		} else if node.right.key > key {
			stacks.Push(node)
			node = node.down
		} else {
			node = node.right
		}
	}
	level := 1
	var download *Node = nil
	for !stacks.Empty() {
		cur := stacks.Peak().(*Node)
		stacks.Pop()
		nodeteam := NewNode(key, val)
		nodeteam.down = download
		download = nodeteam
		if cur.right == nil {
			cur.right = nodeteam
		} else {
			nodeteam.right = cur.right
			cur.right = nodeteam
		}
		if level > MAX_LEVEL {
			break
		}
		var num float32 = rand.Float32()
		if num > 0.5 {
			break
		}
		level++
		if level > s.HighLevel {
			s.HighLevel = level
			highNodedown := NewNode(INT_MIN, 0)
			highNodedown.down = s.HeadNode
			s.HeadNode = highNodedown
			stacks.Push(highNodedown)
		}
	}
}

func (skl *SkipList) Toslice() (score []int, value []interface{}) {
	if skl == nil {
		return nil, nil
	}
	node := skl.HeadNode
	for node.down != nil {
		node = node.down
	}
	node = node.right

	for node != nil {
		score = append(score, node.key)
		value = append(value, node.val)
		node = node.right
	}
	return
}

func (skl *SkipList) GetRange(start int, end int) (score []int, result []interface{}) {
	if skl == nil {
		return nil, nil
	}
	node := skl.HeadNode
	for node.down != nil {
		node = node.down
	}
	node = node.right

	//移到start点；
	wight := start
	for wight > 0 {
		node = node.right
		wight--
	}

	result = []interface{}{}
	score = []int{}
	for i := start; i <= end; i++ {
		score = append(score, node.key)
		result = append(result, node.val)
		node = node.right
	}
	return
}
