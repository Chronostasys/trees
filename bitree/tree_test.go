package bitree

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

type myint int

func (i myint) Hash() int {
	return int(i)
}

func TestTree_BitreeInsert(t *testing.T) {
	tree := &Tree{}
	rands := []int{}
	for i := 19; i >= 0; i -= 1 {
		ran := i
		rands = append(rands, ran)
	}
	rand.Seed(time.Now().UnixMilli())
	rand.Shuffle(len(rands), func(i, j int) {
		rands[i], rands[j] = rands[j], rands[i]
	})
	for _, v := range rands {
		tree.Insert(myint(v))
	}
	sort.Ints(rands)
	i := 0
	tree.Travel(func(val Hasher) {
		if rands[i] != val.Hash() {
			t.Fatalf("wrong travel sequence. expect %d in pos %d, got %d", rands[i], i, val.Hash())
		}
		i++
	})
	tree.Print(true)
}
func TestTree_BitreeDelete(t *testing.T) {
	tree := &Tree{}
	rands := []int{}
	for i := 3; i >= 0; i -= 1 {
		ran := i
		rands = append(rands, ran)
	}
	// rand.Seed(time.Now().UnixMilli())
	// rand.Shuffle(len(rands), func(i, j int) {
	// 	rands[i], rands[j] = rands[j], rands[i]
	// })
	for _, v := range rands {
		tree.Insert(myint(v))
	}
	tree.Delete(1)
	sort.Ints(rands)
	// i := 0
	// tree.Travel(func(val Hasher) {
	// 	if rands[i] != val.Hash() {
	// 		t.Fatalf("wrong travel sequence. expect %d in pos %d, got %d", rands[i], i, val.Hash())
	// 	}
	// 	i++
	// })
	tree.Print(true)
}
