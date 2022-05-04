package skiplist

import (
	"fmt"
	"strconv"
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
	for i := 0; i < 50; i = i + 2 {
		skl.Insert(i, i)
	}
	skl.Delete(12)
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
	// node = skl.HeadNode
	// for node.down != nil {
	// 	node = node.down
	// }
	// node = node.right
	// for node != nil {
	// 	fmt.Println(node.key)
	// 	node = node.right
	// }

	//fmt.Println(skl.Search(11))
}

func TestToSlice(t *testing.T) {
	skl := NewSkipList()
	for i := 0; i < 50; i++ {
		skl.Insert(i, i)
	}
	score, _ := skl.Toslice()
	for k, _ := range score {
		fmt.Println(k)
	}
}

func TestGetRange(t *testing.T) {
	skl := NewSkipList()
	for i := 0; i < 50; i++ {
		skl.Insert(i, byte(i))
	}
	score, val := skl.GetRange(4, 20)
	for _, v := range score {
		fmt.Println(v)
	}
	for _, v := range val {
		fmt.Println(v)
	}

}

func TestGetCount(t *testing.T) {
	skl := NewSkipList()
	for i := 0; i < 50; i++ {
		skl.Insert(i, byte(i))
	}
	count := skl.GetCount(0, 10)
	fmt.Println(count)
}

func TestInsert(t *testing.T) {
	skl := NewSkipList()
	for i := 0; i < 50; i++ {
		skl.Insert(i, byte(i))
	}
	for i := 0; i < 50; i++ {
		skl.Insert(i, byte(i))
	}
	count := skl.GetCount(0, 10)
	fmt.Println(count)
}

func TestFindMember(t *testing.T) {
	skl := NewSkipList()
	for i := 0; i < 50; i++ {
		str := strconv.Itoa(i)
		skl.Insert(i, []byte(str))
	}
	result, ok := skl.FindMember([]byte("15"))
	if ok == false {
		fmt.Println("没存在")
	} else {
		fmt.Println(result)
	}
	result, ok = skl.FindMember([]byte("51"))
	if ok == false {
		fmt.Println("没存在")
	} else {
		fmt.Println(result)
	}

}

func TestRank(t *testing.T) {
	skl := NewSkipList()
	for i := 0; i < 50; i++ {
		str := strconv.Itoa(i)
		skl.Insert(i, []byte(str))
	}
	fmt.Println(skl.GetRank([]byte("100")))
}

func TestDelet(t *testing.T) {
	skl := NewSkipList()
	for i := 0; i < 50; i++ {
		str := strconv.Itoa(i)
		skl.Insert(i, []byte(str))
	}
	for i := 0; i < 50; i++ {
		str := strconv.Itoa(i)
		skl.Insert(i, []byte(str))
	}

}
