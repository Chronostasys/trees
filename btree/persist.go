package btree

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"os"
	"sync"
)

type BinNode struct {
	Childs []int
	Father int
	Right  int
	Vals   []Hasher
}

func init() {
	gob.Register(Int(0))
}

func (n *node) persist(t *Tree) {
	bin := BinNode{
		Childs: make([]int, 0, t.m),
		Vals:   make([]Hasher, 0, t.m-1),
	}
	for _, v := range n.childs {
		bin.Childs = append(bin.Childs, v.fn)
	}
	bin.Vals = n.vals
	if n.father == nil {
		bin.Father = -1
	} else {
		bin.Father = n.father.fn
	}
	if n.right == nil {
		bin.Right = -1
	} else {
		bin.Right = n.right.fn
	}
	if n.f == nil {
		f, _ := os.OpenFile(fmt.Sprintf("%d.idx", n.fn), os.O_CREATE|os.O_RDWR, 0644)
		n.initf(f)
	}
	n.f.Truncate(0)
	n.buf.Reset()
	n.en.Encode(bin)
	if t.takeSnapshot {
		t.snmu.Lock()
		t.snapshot[n.fn] = n.buf.Bytes()
		t.snmu.Unlock()
	}
	n.f.Sync()
}
func (n *node) initf(f *os.File) {
	n.f = f
	n.buf = &bytes.Buffer{}
	n.en = gob.NewEncoder(io.MultiWriter(f, n.buf))
}
func (b *BinNode) loadNode(fn int, t *Tree, loaded map[int]*node) *node {
	n := &node{
		vals:   make([]Hasher, 0, t.m),
		childs: make([]*node, 0),
	}
	for _, v := range b.Childs {
		if loaded[v] == nil {
			f, _ := os.OpenFile(fmt.Sprintf("%d.idx", v), os.O_RDWR, 0644)
			reader := gob.NewDecoder(f)
			bin := &BinNode{}
			reader.Decode(&bin)
			l := bin.loadNode(v, t, loaded)
			n.childs = append(n.childs, l)
			loaded[v] = l
			l.initf(f)
		} else {
			n.childs = append(n.childs, loaded[v])
		}
	}
	n.ensureReversePointer()
	if b.Right == -1 {
		n.right = nil
	} else {
		v := b.Right
		if loaded[v] == nil {
			// panic("xxx")
			f, _ := os.OpenFile(fmt.Sprintf("%d.idx", v), os.O_RDWR, 0644)
			reader := gob.NewDecoder(f)
			bin := &BinNode{}
			reader.Decode(&bin)
			n.right = bin.loadNode(v, t, loaded)
			loaded[v] = n.right
			n.right.initf(f)
		} else {
			n.right = loaded[v]
		}
	}
	n.vals = b.Vals
	n.fn = fn
	return n
}

func LoadSnapshot(sn []byte) *Tree {
	buf := bytes.NewBuffer(sn)
	dec := gob.NewDecoder(buf)
	meta := &TreeMeta{}
	dec.Decode(meta)
	snapshot := map[int][]byte{}
	dec.Decode(&snapshot)
	for k, v := range snapshot {
		go func(k int, v []byte) {
			f, _ := os.OpenFile(fmt.Sprintf("%d.idx", k), os.O_CREATE|os.O_RDWR, 0644)
			f.Write(v)
			f.Sync()
			f.Close()
		}(k, v)
	}
	return loadByMeta(meta)
}

func loadByMeta(meta *TreeMeta) *Tree {
	f1, _ := os.OpenFile(fmt.Sprintf("%d.idx", meta.Rootfn), os.O_RDWR, 0644)
	reader := gob.NewDecoder(f1)
	bin := &BinNode{}
	reader.Decode(&bin)
	t := &Tree{
		m:        meta.M,
		total:    meta.Total,
		gfn:      meta.Gfn,
		edge:     int(math.Ceil(float64((meta.M-1))/2)) - 1,
		fs:       map[*node]struct{}{},
		snapshot: map[int][]byte{},
		snmu:     &sync.Mutex{},
	}
	loaded := map[int]*node{}
	t.root = bin.loadNode(meta.Rootfn, t, loaded)
	t.root.initf(f1)
	t.first = loaded[meta.First]
	return t
}
func Load() *Tree {
	meta := &TreeMeta{}
	f, _ := os.OpenFile(".meta", os.O_CREATE|os.O_RDWR, 0644)
	enc := gob.NewDecoder(f)
	enc.Decode(meta)
	return loadByMeta(meta)
}
func (t *Tree) Persist() {
	wg := sync.WaitGroup{}
	le := len(t.fs)
	wg.Add(le)
	for n := range t.fs {
		go func(n *node) {
			n.persist(t)
			wg.Done()
		}(n)
		delete(t.fs, n)
	}
	meta := TreeMeta{
		Rootfn: t.root.fn,
		Total:  t.total,
		M:      t.m,
		Gfn:    t.gfn,
		First:  t.first.fn,
	}
	if t.f == nil {
		t.f, _ = os.OpenFile(".meta", os.O_CREATE|os.O_RDWR, 0644)
		t.buf = &bytes.Buffer{}
		t.en = gob.NewEncoder(io.MultiWriter(t.f, t.buf))
	}
	t.buf.Reset()
	t.en.Encode(meta)
	t.f.Sync()
	wg.Wait()
}
func (t *Tree) PersistWithSnapshot() []byte {
	t.takeSnapshot = true
	t.Persist()
	buf := &bytes.Buffer{}
	buf.Write(t.buf.Bytes())
	enc := gob.NewEncoder(buf)
	enc.Encode(t.snapshot)
	return buf.Bytes()
}

type TreeMeta struct {
	Rootfn int
	Total  int
	M      int
	Gfn    int
	First  int
}
