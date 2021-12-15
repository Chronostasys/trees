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
	tree := Make(256)
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
	tree.Iterate(func(val Item) {
		if rands[i] != val.(Int).Int() {
			t.Fatalf("wrong travel sequence. expect %d in pos %d, got %d", rands[i], i, val.(Int).Int())
		}
		i++
	})
	if i != len(rands) {
		t.Fatalf("only travel %d items, expected %d items", i, len(rands))
	}
	// tree.Print(true)
}

func BenchmarkInsert(b *testing.B) {
	tree := Make(256)
	arr := rand.Perm(b.N)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.Insert(Int(arr[n]))
	}
}
func BenchmarkGoogleInsert(b *testing.B) {
	tree := btree.New(256)
	arr := rand.Perm(b.N)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.ReplaceOrInsert(btree.Int(arr[n]))
	}
}

func TestTree_BtreeDelete(t *testing.T) {
	tree := Make(256)
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
		tree.Delete(Int(rands[i]))
	}
	rands = rands[500000:]
	sort.Ints(rands)
	i := 0
	tree.Iterate(func(val Item) {
		if rands[i] != val.(Int).Int() {
			t.Fatalf("wrong travel sequence. expect %d in pos %d, got %d", rands[i], i, val.(Int).Int())
		}
		i++
	})
	// tree.Print(true)
}

func BenchmarkDelete(b *testing.B) {
	tree := Make(256)
	arr := rand.Perm(b.N)
	for n := 0; n < b.N; n++ {
		tree.Insert(Int(arr[n]))
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.Delete(Int(arr[n]))
	}
}
func BenchmarkGoogleDelete(b *testing.B) {
	tree := btree.New(256)
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
	tree := Make(256)
	arr := rand.Perm(b.N)
	for n := 0; n < b.N; n++ {
		tree.Insert(Int(arr[n]))
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.Search(Int(arr[n]))
	}
}
func BenchmarkGoogleSearch(b *testing.B) {
	tree := btree.New(256)
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
	tree := Make(256)
	arr := rand.Perm(1000000)
	for i := 0; i < 1000000; i++ {
		tree.Insert(Int(arr[i]))
	}
	for _, v := range arr {
		val := tree.Search(Int(v))
		if val.(Int) != Int(v) {
			t.Fatal("search get wrong value. Expect", v, "got", val.(Int))
		}
	}
	for i := 0; i < 500000; i++ {
		tree.Delete(Int(arr[i]))
	}
	for _, v := range arr[500000:] {
		val := tree.Search(Int(v))
		if val.(Int) != Int(v) {
			t.Fatal("search get wrong value. Expect", v, "got", val.(Int))
		}
	}
	for _, v := range arr[:500000] {
		val := tree.Search(Int(v))
		if val != nil {
			t.Fatal("search get wrong value. Expect nil", "got", val.(Int))
		}
	}
	// tree.Print(true)
}

func Test_Persist(t *testing.T) {
	tree := MakePersist(256)
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
	sn := tree.PersistWithSnapshot("test/t-")
	tree = LoadSnapshot(sn, "test/t-")
	tree.Iterate(func(val Item) {
		if rands[i] != val.(Int).Int() {
			t.Fatalf("wrong travel sequence. expect %d in pos %d, got %d", rands[i], i, val.(Int).Int())
		}
		i++
	})
}

func TestTree_LargerOrEq(t *testing.T) {
	tree := Make(256)
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
	n := 1000
	tree.LargerOrEq(Int(n), 2000, func(i Item) bool {
		if n != i.(Int).Int() {
			t.Fatalf("expect %d, got %d", n, i.(Int).Int())
		}
		n++
		return true
	})
	if n != 3000 {
		t.Errorf("expect n=3000 after test, got n=%d", n)
	}
	n = 999
	tree.Larger(Int(n), 2000, func(i Item) bool {
		if n+1 != i.(Int).Int() {
			t.Fatalf("expect %d, got %d", n+1, i.(Int).Int())
		}
		n++
		return true
	})
	if n != 3000-1 {
		t.Fatalf("expect n=2999 after test, got n=%d", n)
	}
}
