package btree

import (
	"math"
	"os"
	"sort"
	"sync"
)

func Make(m int) *Tree {
	return &Tree{
		m:    m,
		edge: int(math.Ceil(float64((m-1))/2)) - 1,
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
		vals:   make([]Item, 0, m),
		childs: make([]*node, 0), // all leaves do not have childs. So init it to zero minimize allocation. (if you meant to set it's len, set m+1)
		fn:     t.gfn,
	}
}

func (t *Tree) Insert(val Item) {
	if t.root == nil {
		t.root = t.makeBNode()
		t.root.vals = append(t.root.vals, val)
		t.first = t.root
		t.total++
	} else {
		t.root.insert(t, val)
	}
}

func itemEQ(i1, i2 Item) bool {
	return i1.EQ(i2)
}

func (n *node) insert(t *Tree, val Item) {
	if len(n.childs) == 0 {
		index := n.biSearch(val)
		// update situation
		if index-1 > -1 && itemEQ(n.vals[index-1], val) {
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
			if len(nchilds) != 0 { // 向上分裂
				lf.childs = nchilds[:t.m/2+1]
				ri.childs = append(ri.childs, nchilds[t.m/2+1:]...)
				ri.ensureReversePointer()
			} else {
				ri.right = n.right
				lf.right = ri
			}
			if father != nil {
				ri.father = father
				idx := father.biSearch(nvals[0])
				last := len(father.vals) - 1
				father.vals = append(father.vals, father.vals[last])
				ed := father.vals[idx:]
				copy(father.vals[idx+1:], ed)
				father.vals[idx] = ri.vals[0].Key()
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
				t.root.vals = append(t.root.vals, nvals[t.m/2].Key())
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
				n.vals = append(n.vals[:0], ri.vals[0].Key())
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
	idx := n.biSearch(val)
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
func (t *Tree) Iterate(job func(val Item)) {
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
func (n *node) biSearch(item Item) int {
	if n == nil || n.vals == nil {
		return -1
	}
	return sort.Search(len(n.vals), func(i int) bool {
		return item.Less(n.vals[i])
	})
}

func (t *Tree) Delete(item Item) {
	if t.root == nil {
		return
	}
	t.root.delete(t, item)
}

func (n *node) delete(t *Tree, item Item) {
	// leaf node
	if len(n.childs) == 0 {
		index := n.biSearch(item) - 1
		// exist
		if index != -1 && itemEQ(n.vals[index], item) {
			first := n.vals[0]
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
						n.vals = append(n.vals, nil)
						copy(n.vals[1:], n.vals[:len(n.vals)-1])
						n.vals[0] = father.vals[idx-1]
						father.vals[idx-1] = bro.vals[last].Key()
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
					father.vals[idx-1] = bro.vals[last].Key()
					lenn := len(n.vals)
					n.vals = append(n.vals, nil)
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
						father.vals[idx] = bro.vals[0].Key()
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
					father.vals[idx] = bro.vals[0].Key()
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
	idx := n.biSearch(item)
	if len(n.childs) <= idx {
		// not exist
		return
	}
	n.childs[idx].delete(t, item)

}

func (t *Tree) Search(item Item) Item {
	if t.root == nil {
		return nil
	}
	return t.root.search(item)
}
func (n *node) search(item Item) Item {
	idx := n.biSearch(item)
	if len(n.childs) == 0 {
		if idx-1 < 0 || !itemEQ(n.vals[idx-1], item) {
			return nil
		}
		return n.vals[idx-1]
	}
	return n.childs[idx].search(item)
}
func (t *Tree) Len() int {
	return t.total
}

func (t *Tree) Larger(item Item, max int, callback func(Item)) {
	if t.root == nil {
		return
	}
	t.root.largerOrEq(item, max, callback, false)
}
func (n *node) largerOrEq(item Item, max int, callback func(Item), eq bool) {
	idx := n.biSearch(item)
	if len(n.childs) == 0 {
		start := idx - 1
		if !eq && n.vals[idx-1].EQ(item) {
			start = idx
			max = max + 1
		}
		ri := idx - 1 + max
		if idx-1+max > len(n.vals) {
			if start != len(n.vals) {
				for _, v := range n.vals[start:] {
					callback(v)
				}
			}
			nx := ri
			for {
				nx = nx - len(n.vals)
				if n.right == nil {
					println()
				}
				n = n.right
				if n == nil || nx <= 0 {
					return
				}
				for i, v := range n.vals {
					if i < nx {
						callback(v)
					}
				}
			}
		} else {
			for _, v := range n.vals[idx-1 : ri] {
				callback(v)

			}
			return
		}
	}
	n.childs[idx].largerOrEq(item, max, callback, eq)
}
func (t *Tree) LargerOrEq(item Item, max int, callback func(Item)) {
	if t.root == nil {
		return
	}
	t.root.largerOrEq(item, max, callback, true)
}
