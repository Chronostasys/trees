package bitree

import (
	"sort"
	"testing"
)

type myint int

func (i myint) Hash() int {
	return int(i)
}

func TestTree_Bitree(t *testing.T) {
	tree := &Tree{}
	rands := []int{}
	for i := 6; i > 1; i-- {
		ran := i
		rands = append(rands, ran)
		tree.Insert(myint(ran))
	}
	sort.Ints(rands)
	i := 0
	tree.Travel(func(val Hasher) {
		if rands[i] != val.Hash() {
			t.Fatalf("wrong travel sequence. expect %d in pos %d, got %d", rands[i], i, val.Hash())
		}
		i++
	})
	tree.Print()
}
