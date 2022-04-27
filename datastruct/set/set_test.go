package set

import (
	"fmt"
	"testing"
)

func TestNewDict(t *testing.T) {
	set := NewSet()
	set.Add("123")
	state := set.Has("123")
	if state == false {
		t.Error("set add error")
	}
	set.Add("323")
	if set.Len() != 2 {
		t.Error("set len error")
	}

}

func TestUnion(t *testing.T) {
	slice := make([]string, 3)
	slice[0] = "name"
	slice[1] = "age"
	slice[2] = "class"
	set1 := NewSetWithVal(slice...)
	set2 := NewSet()
	set2.Add("name")
	set2.Add("grade")
	set3 := set1.Union(set2)
	keys := set3.dict.Keys()
	fmt.Println(len(keys))
	for _, v := range keys {
		fmt.Println(v)
	}
}
