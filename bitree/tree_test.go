package bitree

import (
	"math/rand"
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
	for i := 0; i < 100; i++ {
		ran := rand.Int()
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
}
