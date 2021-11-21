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
		tree.Insert(Int(v))
	}
	sort.Ints(rands)
	i := 0
	tree.Iterate(func(val Hasher) {
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
		tree.Insert(Int(arr[n]))
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
		tree.Insert(Int(v))
	}
	rand.Seed(time.Now().UnixMilli())
	rand.Shuffle(len(rands), func(i, j int) {
		rands[i], rands[j] = rands[j], rands[i]
	})
	for i := 0; i < 500000; i++ {
		tree.Delete(rands[i])
	}
	rands = rands[500000:]
	sort.Ints(rands)
	i := 0
	tree.Iterate(func(val Hasher) {
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
		tree.Insert(Int(arr[n]))
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
		tree.Insert(Int(arr[n]))
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
	tree := Make(1024)
	arr := rand.Perm(1000000)
	for i := 0; i < 1000000; i++ {
		tree.Insert(Int(arr[i]))
	}
	for _, v := range arr {
		val := tree.Search(v)
		if val.(Int) != Int(v) {
			t.Fatal("search get wrong value. Expect", v, "got", val.(Int))
		}
	}
	for i := 0; i < 500000; i++ {
		tree.Delete(arr[i])
	}
	for _, v := range arr[500000:] {
		val := tree.Search(v)
		if val.(Int) != Int(v) {
			t.Fatal("search get wrong value. Expect", v, "got", val.(Int))
		}
	}
	for _, v := range arr[:500000] {
		val := tree.Search(v)
		if val != nil {
			t.Fatal("search get wrong value. Expect nil", "got", val.(Int))
		}
	}
	// tree.Print(true)
}

func Test_Persist(t *testing.T) {
	tree := MakePersist(1024)
	rands := []int{}
	for i := 100000; i >= 0; i -= 1 {
		ran := i
		rands = append(rands, ran)
	}
	rand.Seed(time.Now().UnixMilli())
	rand.Shuffle(len(rands), func(i, j int) {
		rands[i], rands[j] = rands[j], rands[i]
	})
	for _, v := range rands {
		tree.Insert(Int(v))
	}
	sort.Ints(rands)
	i := 0
	t1 := time.Now()
	sn := tree.PersistWithSnapshot("test")
	l := len(sn)
	println(l)
	println(time.Since(t1).String())
	t1 = time.Now()
	tree = LoadSnapshot(sn, "test")
	println(time.Since(t1).String())
	tree.Iterate(func(val Hasher) {
		if rands[i] != val.Hash() {
			t.Fatalf("wrong travel sequence. expect %d in pos %d, got %d", rands[i], i, val.Hash())
		}
		i++
	})
}
