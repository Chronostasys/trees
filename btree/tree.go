package btree

import (
	"bytes"
	"encoding/gob"
	"math"
	"os"
	"sort"
	"sync"
)

type Int int

func (i Int) Hash() int {
	return int(i)
}

type node struct {
	childs []*node
	father *node
	vals   []Hasher
	right  *node
	fn     int
	f      *os.File
	en     *gob.Encoder
	buf    *bytes.Buffer
}
type Hasher interface {
	Hash() int
}

type Tree struct {
	snmu         *sync.Mutex
	root         *node
	total        int
	m            int // 阶
	edge         int
	first        *node
	gfn          int
	fs           map[*node]struct{}
	f            *os.File
	buf          *bytes.Buffer
	en           *gob.Encoder
	takeSnapshot bool
	snapshot     map[int][]byte
	persist      bool
}

func Make(m int) *Tree {
	return &Tree{
		m:        m,
		edge:     int(math.Ceil(float64((m-1))/2)) - 1,
		fs:       map[*node]struct{}{},
		snapshot: map[int][]byte{},
		snmu:     &sync.Mutex{},
	}
}
func MakePersist(m int) *Tree {
	return &Tree{
		m:        m,
		edge:     int(math.Ceil(float64((m-1))/2)) - 1,
		fs:       map[*node]struct{}{},
		snapshot: map[int][]byte{},
		snmu:     &sync.Mutex{},
		persist:  true,
	}
}

func (t *Tree) makeBNode() *node {
	m := t.m
	defer func() {
		t.gfn++
	}()
	return &node{
		vals:   make([]Hasher, 0, m),
		childs: make([]*node, 0), // all leaves do not have childs. So init it to zero minimize allocation. (if you meant to set it's len, set m+1)
		fn:     t.gfn,
	}
}

func (t *Tree) Insert(val Hasher) {
	if t.root == nil {
		t.root = t.makeBNode()
		t.root.vals = append(t.root.vals, val)
		t.first = t.root
		t.total++
	} else {
		t.root.insert(t, val)
	}
}

func (n *node) insert(t *Tree, val Hasher) {
	if len(n.childs) == 0 {
		index := n.biSearch(val.Hash())
		// update situation
		if index-1 > -1 && n.vals[index-1].Hash() == val.Hash() {
			n.vals[index-1] = val
			if t.persist {
				t.fs[n] = struct{}{}
			}
			return
		}
		last := len(n.vals) - 1
		if index == len(n.vals) {
			n.vals = append(n.vals, val)
		} else {
			n.vals = append(n.vals, n.vals[last])
			copy(n.vals[index+1:], n.vals[index:last])
			n.vals[index] = val
		}
		if t.persist {
			t.fs[n] = struct{}{}
		}
		t.total++
	START:
		if len(n.vals) == t.m {
			father := n.father
			nvals := n.vals[:]
			nchilds := n.childs[:]
			lf := n
			lf.vals = nvals[:t.m/2]
			ri := t.makeBNode()
			ri.vals = append(ri.vals, nvals[t.m/2:]...)
			lf.right = ri
			if len(nchilds) != 0 { // 向上分裂
				lf.childs = nchilds[:t.m/2+1]
				ri.childs = append(ri.childs, nchilds[t.m/2+1:]...)
				ri.ensureReversePointer()
			}
			if father != nil {
				ri.father = father
				idx := father.biSearch(nvals[0].Hash())
				last := len(father.vals) - 1
				father.vals = append(father.vals, father.vals[last])
				ed := father.vals[idx:]
				copy(father.vals[idx+1:], ed)
				father.vals[idx] = Int(ri.vals[0].Hash())
				last = len(father.childs) - 1
				father.childs = append(father.childs, father.childs[last])
				copy(father.childs[idx+2:], father.childs[idx+1:])
				father.childs[idx+1] = ri
				if len(nchilds) != 0 { // 向上分裂
					copy(ri.vals[:len(ri.vals)-1], ri.vals[1:])
					ri.vals = ri.vals[:len(ri.vals)-1]
				}
				if t.persist {
					t.fs[lf] = struct{}{}
					t.fs[ri] = struct{}{}
					t.fs[father] = struct{}{}
				}
				n = father
				goto START
			} else if len(nchilds) != 0 { // 向上分裂
				copy(ri.vals[:len(ri.vals)-1], ri.vals[1:])
				ri.vals = ri.vals[:len(ri.vals)-1]
				t.root = t.makeBNode()
				t.root.vals = append(t.root.vals, Int(nvals[t.m/2].Hash()))
				t.root.childs = append(t.root.childs, lf, ri)
				lf.father = t.root
				ri.father = t.root
				if t.persist {
					t.fs[t.root] = struct{}{}
					t.fs[lf] = struct{}{}
					t.fs[ri] = struct{}{}
				}
			} else {
				n = t.makeBNode()
				n.vals = append(n.vals[:0], Int(ri.vals[0].Hash()))
				n.childs = append(n.childs, lf, ri)
				t.root = n
				lf.father = n
				ri.father = n
				t.first = lf

				if t.persist {
					t.fs[t.root] = struct{}{}
					t.fs[lf] = struct{}{}
					t.fs[ri] = struct{}{}
				}
			}
		}
		return
	}
	idx := n.biSearch(val.Hash())
	if len(n.childs) <= idx {
		no := t.makeBNode()
		no.father = n
		n.childs = append(n.childs, no)
	}
	n.childs[idx].insert(t, val)
}
func (n *node) ensureReversePointer() {
	for _, v := range n.childs {
		v.father = n
	}
}
func (t *Tree) Iterate(job func(val Hasher)) {
	if t.root == nil {
		return
	}
	e := t.first
	for {
		if e == nil {
			break
		}
		for _, v := range e.vals {
			job(v)
		}
		e = e.right
	}
}

// binary search.
func (n *node) biSearch(hash int) int {
	if n == nil || n.vals == nil {
		return -1
	}
	return sort.Search(len(n.vals), func(i int) bool {
		return n.vals[i].Hash() > hash
	})
}

func (t *Tree) Delete(hash int) {
	if t.root == nil {
		return
	}
	t.root.delete(t, hash)
}

func (n *node) delete(t *Tree, hash int) {
	// leaf node
	if len(n.childs) == 0 {
		index := n.biSearch(hash) - 1
		// exist
		if index != len(n.vals) && n.vals[index].Hash() == hash {
			first := n.vals[0].Hash()
			if index == len(n.vals)-1 {
				n.vals = n.vals[:index]
			} else {
				n.vals = append(n.vals[:index], n.vals[index+1:]...)
			}
			t.total--
			if t.persist {
				t.fs[n] = struct{}{}
			}
		START:
			if len(n.vals) >= t.edge {
				// valid leaf, return directly
				return
			}

			// node try to borrow val from brother
			father := n.father
			// root node, return
			if father == nil {
				return
			}
			var bro *node
			idx := father.biSearch(first)
			left := true
			if idx-1 > -1 {
				bro = father.childs[idx-1]
				if len(bro.vals) > t.edge {
					// can borrow
					last := len(bro.vals) - 1
					if len(n.childs) > 0 {
						// index nodes
						n.vals = append(n.vals, Int(0))
						copy(n.vals[1:], n.vals[:len(n.vals)-1])
						n.vals[0] = father.vals[idx-1]
						father.vals[idx-1] = Int(bro.vals[last].Hash())
						bro.vals = bro.vals[:last]
						n.childs = append(n.childs, nil)
						copy(n.childs[1:], n.childs[:len(n.childs)])
						n.childs[0] = bro.childs[len(bro.childs)-1]
						n.childs[0].father = n
						bro.childs = bro.childs[:len(bro.childs)-1]
						if t.persist {
							t.fs[n] = struct{}{}
							t.fs[bro] = struct{}{}
							t.fs[father] = struct{}{}
						}
						return
					}
					father.vals[idx-1] = Int(bro.vals[last].Hash())
					lenn := len(n.vals)
					n.vals = append(n.vals, Int(0))
					copy(n.vals[1:], n.vals[:lenn])
					n.vals[0] = bro.vals[last]
					bro.vals = bro.vals[:last]
					if t.persist {
						t.fs[n] = struct{}{}
						t.fs[bro] = struct{}{}
						t.fs[father] = struct{}{}
					}
					return
				}
			}
			if idx+1 < len(father.childs) {
				left = false
				bro = father.childs[idx+1]
				if len(bro.vals) > t.edge {
					// can borrow
					if len(n.childs) > 0 {
						// index nodes
						n.vals = append(n.vals, father.vals[idx])
						father.vals[idx] = Int(bro.vals[0].Hash())
						copy(bro.vals[:len(bro.vals)-1], bro.vals[1:])
						bro.vals = bro.vals[:len(bro.vals)-1]
						n.childs = append(n.childs, bro.childs[0])
						bro.childs[0].father = n
						copy(bro.childs[:len(bro.childs)-1], bro.childs[1:])
						bro.childs = bro.childs[:len(bro.childs)-1]
						if t.persist {
							t.fs[n] = struct{}{}
							t.fs[bro] = struct{}{}
							t.fs[father] = struct{}{}
						}
						return
					}
					n.vals = append(n.vals, bro.vals[0])
					copy(bro.vals[:len(bro.vals)-1], bro.vals[1:])
					bro.vals = bro.vals[:len(bro.vals)-1]
					father.vals[idx] = Int(bro.vals[0].Hash())
					if t.persist {
						t.fs[n] = struct{}{}
						t.fs[bro] = struct{}{}
						t.fs[father] = struct{}{}
					}
					return
				}
			}
			// failed to borrow, merge it!
			if bro == nil {
				// seems it's the root, check it
				if n != t.root {
					panic("not root!")
				}
				return
			}
			if left {
				if len(n.childs) > 0 {
					// index merge
					bro.vals = append(bro.vals, father.vals[idx-1])
					bro.childs = append(bro.childs, n.childs...)
					for _, v := range n.childs {
						v.father = bro
					}
				}
				bro.vals = append(bro.vals, n.vals...)
				father.vals = append(father.vals[:idx-1], father.vals[idx:]...)
				father.childs = append(father.childs[:idx], father.childs[idx+1:]...)
				bro.right = n.right
				if t.persist {
					t.fs[bro] = struct{}{}
					t.fs[father] = struct{}{}
					delete(t.fs, n)
					go func(n *node) {
						name := n.f.Name()
						n.f.Close()
						os.Remove(name)
					}(n)
				}

				n = father
				if t.root == n {
					if len(n.childs) == 1 {
						t.root = n.childs[0]
						t.root.father = nil
						if t.persist {
							t.fs[t.root] = struct{}{}
						}
					}
					return
				}

				goto START
			} else {
				if len(n.childs) > 0 {
					// index merge
					n.vals = append(n.vals, father.vals[idx])
					n.childs = append(n.childs, bro.childs...)
					for _, v := range bro.childs {
						v.father = n
					}
					n.right = bro.right
				}
				n.vals = append(n.vals, bro.vals...)
				father.vals = append(father.vals[:idx], father.vals[idx+1:]...)
				father.childs = append(father.childs[:idx+1], father.childs[idx+2:]...)
				n.right = bro.right
				if t.persist {
					t.fs[n] = struct{}{}
					t.fs[father] = struct{}{}
					delete(t.fs, bro)
					go func(n *node) {
						name := n.f.Name()
						n.f.Close()
						os.Remove(name)
					}(bro)
				}
				n = father
				if t.root == n {
					if len(n.childs) == 1 {
						t.root = n.childs[0]
						t.root.father = nil
						if t.persist {
							t.fs[t.root] = struct{}{}
						}
					}
					return
				}
				goto START
			}
		}
		return
	}
	// index node
	idx := n.biSearch(hash)
	if len(n.childs) <= idx {
		// not exist
		return
	}
	n.childs[idx].delete(t, hash)

}

func (t *Tree) Search(hash int) Hasher {
	if t.root == nil {
		return nil
	}
	return t.root.search(hash)
}
func (n *node) search(hash int) Hasher {
	idx := n.biSearch(hash)
	if len(n.childs) == 0 {
		if idx-1 < 0 || n.vals[idx-1].Hash() != hash {
			return nil
		}
		return n.vals[idx-1]
	}
	return n.childs[idx].search(hash)
}
