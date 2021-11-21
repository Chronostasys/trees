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
	Childs []int64
	Father int64
	Right  int64
	Vals   []Hasher
}

func init() {
	gob.Register(Int(0))
}

func (n *node) persist(t *Tree, prefix string) {
	bin := BinNode{
		Childs: make([]int64, 0, t.m),
		Vals:   make([]Hasher, 0, t.m-1),
	}
	for _, v := range n.childs {
		bin.Childs = append(bin.Childs, int64(v.fn))
	}
	bin.Vals = n.vals
	if n.father == nil {
		bin.Father = -1
	} else {
		bin.Father = int64(n.father.fn)
	}
	if n.right == nil {
		bin.Right = -1
	} else {
		bin.Right = int64(n.right.fn)
	}
	if n.f == nil {
		f, _ := os.OpenFile(fmt.Sprintf("%s%d.idx", prefix, n.fn), os.O_CREATE|os.O_RDWR, 0644)
		n.initf(f)
	}
	n.f.Truncate(0)
	n.buf.Reset()
	n.en = gob.NewEncoder(io.MultiWriter(n.f, n.buf))
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
}
func (b *BinNode) loadNode(fn int, t *Tree, loaded map[int]*node, prefix string) *node {
	n := &node{
		vals:   make([]Hasher, 0, t.m),
		childs: make([]*node, 0),
	}
	for _, u := range b.Childs {
		v := int(u)
		if loaded[v] == nil {
			f, _ := os.OpenFile(fmt.Sprintf("%s%d.idx", prefix, v), os.O_RDWR, 0644)
			reader := gob.NewDecoder(f)
			bin := &BinNode{}
			reader.Decode(&bin)
			l := bin.loadNode(v, t, loaded, prefix)
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
		v := int(b.Right)
		if loaded[v] == nil {
			// panic("xxx")
			f, _ := os.OpenFile(fmt.Sprintf("%s%d.idx", prefix, v), os.O_RDWR, 0644)
			reader := gob.NewDecoder(f)
			bin := &BinNode{}
			reader.Decode(&bin)
			n.right = bin.loadNode(v, t, loaded, prefix)
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

func LoadSnapshot(sn []byte, prefix string) *Tree {
	buf := bytes.NewBuffer(sn)
	dec := gob.NewDecoder(buf)
	meta := &TreeMeta{}
	err := dec.Decode(meta)
	if err != nil {
		return nil
	}
	snapshot := map[int][]byte{}
	dec.Decode(&snapshot)
	for k, v := range snapshot {
		go func(k int, v []byte) {
			f, _ := os.OpenFile(fmt.Sprintf("%s%d.idx", prefix, k), os.O_CREATE|os.O_RDWR, 0644)
			f.Write(v)
			f.Sync()
			f.Close()
		}(k, v)
	}
	return loadByMeta(meta, prefix)
}

func loadByMeta(meta *TreeMeta, prefix string) *Tree {
	f1, err := os.OpenFile(fmt.Sprintf("%s%d.idx", prefix, meta.Rootfn), os.O_RDWR, 0644)
	if err != nil {
		return nil
	}
	reader := gob.NewDecoder(f1)
	bin := &BinNode{}
	err = reader.Decode(&bin)
	if err != nil {
		return nil
	}
	t := &Tree{
		m:        meta.M,
		total:    meta.Total,
		gfn:      meta.Gfn,
		edge:     int(math.Ceil(float64((meta.M-1))/2)) - 1,
		fs:       map[*node]struct{}{},
		snapshot: map[int][]byte{},
		snmu:     &sync.Mutex{},
		persist:  true,
	}
	loaded := map[int]*node{}
	t.root = bin.loadNode(meta.Rootfn, t, loaded, prefix)
	t.root.initf(f1)
	loaded[t.root.fn] = t.root
	t.first = loaded[meta.First]
	return t
}
func Load(prefix string) *Tree {
	meta := &TreeMeta{}
	f, _ := os.OpenFile(prefix+".meta", os.O_CREATE|os.O_RDWR, 0644)
	enc := gob.NewDecoder(f)
	err := enc.Decode(meta)
	if err != nil {
		return nil
	}
	return loadByMeta(meta, prefix)
}
func (t *Tree) Persist(prefix string) {
	wg := sync.WaitGroup{}
	le := len(t.fs)
	wg.Add(le)
	for n := range t.fs {
		go func(n *node) {
			n.persist(t, prefix)
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
		t.f, _ = os.OpenFile(prefix+".meta", os.O_CREATE|os.O_RDWR, 0644)
		t.buf = &bytes.Buffer{}
		// t.en = gob.NewEncoder(io.MultiWriter(t.f, t.buf))
	}
	t.buf.Reset()
	t.en = gob.NewEncoder(io.MultiWriter(t.f, t.buf))
	t.en.Encode(meta)
	t.f.Sync()
	wg.Wait()
}
func (t *Tree) PersistWithSnapshot(prefix string) []byte {
	t.takeSnapshot = true
	t.Persist(prefix)
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
