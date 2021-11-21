package btree

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/google/btree"
)

func init() {
	rand.Seed(time.Now().UnixMilli())
}

func TestTree_BtreeInsert(t *testing.T) {
	tree := Make(1024)
	rands := []int{}
	for i := 1000000; i >= 0; i -= 1 {
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
	tree := Make(1024)
	arr := rand.Perm(b.N)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.Insert(myint(arr[n]))
	}
}
func BenchmarkGoogleInsert(b *testing.B) {
	tree := btree.New(1024)
	arr := rand.Perm(b.N)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.ReplaceOrInsert(btree.Int(arr[n]))
	}
}

func TestTree_BtreeDelete(t *testing.T) {
	tree := Make(1024)
	rands := []int{}
	for i := 1000000; i >= 0; i -= 1 {
		ran := i
		rands = append(rands, ran)
	}
	for _, v := range rands {
		tree.Insert(myint(v))
	}
	rand.Seed(time.Now().UnixMilli())
	rand.Shuffle(len(rands), func(i, j int) {
		rands[i], rands[j] = rands[j], rands[i]
	})
	for i := 0; i < 50000; i++ {
		tree.Delete(rands[i])
	}
	rands = rands[50000:]
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

func BenchmarkDelete(b *testing.B) {
	tree := Make(1024)
	arr := rand.Perm(b.N)
	for n := 0; n < b.N; n++ {
		tree.Insert(myint(arr[n]))
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.Delete(arr[n])
	}
}
func BenchmarkGoogleDelete(b *testing.B) {
	tree := btree.New(1024)
	arr := rand.Perm(b.N)
	for n := 0; n < b.N; n++ {
		tree.ReplaceOrInsert(btree.Int(arr[n]))
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.Delete(btree.Int(arr[n]))
	}
}

func BenchmarkSearch(b *testing.B) {
	tree := Make(1024)
	arr := rand.Perm(b.N)
	for n := 0; n < b.N; n++ {
		tree.Insert(myint(arr[n]))
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.Search(arr[n])
	}
}
func BenchmarkGoogleSearch(b *testing.B) {
	tree := btree.New(1024)
	arr := rand.Perm(b.N)
	for n := 0; n < b.N; n++ {
		tree.ReplaceOrInsert(btree.Int(arr[n]))
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.Get(btree.Int(arr[n]))
	}
}

func TestTree_BtreeSearch(t *testing.T) {
	tree := Make(30)
	for i := 1000; i >= 0; i -= 1 {
		tree.Insert(myint(i))
	}
	val := tree.Search(500)
	if val.(myint) != 500 {
		t.Fatal("search get wrong value. Expect", 500, "got", val.(myint))
	}
	val = tree.Search(-1)
	if val != nil {
		t.Fatal("search get wrong value. Expect nil", "got", val.(myint))
	}
	// tree.Print(true)
}
