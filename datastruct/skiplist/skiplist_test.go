package skiplist

import (
	"fmt"
	"testing"
)

func TestNewSkipList(t *testing.T) {
	// skl := NewSkipList()
	// skl.Insert(1, 1)
	// skl.Insert(2, 2)
	// skl.Insert(3, 3)
	// skl.Insert(4, 4)
	// for i := 0; i <= 4; i++ {
	// 	fmt.Println(skl.Search(i))
	// }
}

func TestLevel(t *testing.T) {
	skl := NewSkipList()
	for i := 0; i < 50; i = i + 2 {
		skl.Insert(i, i)
	}
	node := skl.HeadNode
	for node != nil {
		cur := node
		for cur != nil {
			fmt.Print(cur.key)
			fmt.Print("->")
			cur = cur.right
		}
		fmt.Println(" ")
		node = node.down
	}

	fmt.Println(skl.Search(11))
}
