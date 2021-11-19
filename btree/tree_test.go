package btree

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestTree_BtreeInsert(t *testing.T) {
	tree := Make(3)
	rands := []int{}
	for i := 60; i >= 0; i -= 1 {
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
	tree.Travel(func(val Hasher, level int) {
		if rands[i] != val.Hash() {
			t.Fatalf("wrong travel sequence. expect %d in pos %d, got %d", rands[i], i, val.Hash())
		}
		i++
	})
	// tree.Print(true)
}

func BenchmarkInsert(b *testing.B) {
	tree := Make(1000)
	for n := 0; n < b.N; n++ {
		tree.Insert(myint(n))
	}
}

func BenchmarkMap(b *testing.B) {
	m := make(map[myint]struct{})
	s := struct{}{}
	for n := 0; n < b.N; n++ {
		m[myint(n)] = s
	}
}
