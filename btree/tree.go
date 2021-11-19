package btree

type myint int

func (i myint) Hash() int {
	return int(i)
}

type node struct {
	childs []*node
	father *node
	vals   []Hasher
}

type Hasher interface {
	Hash() int
}

type Tree struct {
	root  *node
	total int
	m     int // 阶
}

func Make(m int) *Tree {
	return &Tree{m: m}
}

func (t *Tree) Insert(val Hasher) {
	if t.root == nil {
		t.root = &node{
			vals: []Hasher{val},
		}
		t.total++
	} else {
		t.root.insert(t, val)
	}
}

func (n *node) insert(t *Tree, val Hasher) {
	if len(n.childs) == 0 {
		vals := []Hasher{}
		for i, v := range n.vals {
			if v.Hash() == val.Hash() {
				n.vals[i] = val
				return
			} else if v.Hash() > val.Hash() {
				vals = append(vals, val)
				t.total++
				vals = append(vals, n.vals[i:]...)
				break
			}
			vals = append(vals, v)
		}
		if len(vals) == len(n.vals) {
			vals = append(vals, val)
			t.total++
		}
		n.vals = vals
	START:
		if len(n.vals) == t.m {
			lf := &node{
				vals:   n.vals[:t.m/2],
				father: n,
			}
			ri := &node{
				vals:   n.vals[t.m/2:],
				father: n,
			}
			if len(n.childs) != 0 { // 向上分裂
				lf.childs = n.childs[:t.m/2+1]
				ri.childs = n.childs[t.m/2+1:]
				lf.ensureReversePointer()
				ri.ensureReversePointer()
			}
			if n.father != nil {
				idx := -1
				father := n.father
				lf.father = father
				ri.father = n.father
				for i, v := range father.vals {
					if v.Hash() > n.vals[0].Hash() {
						idx = i
						break
					}
				}
				if idx == -1 {
					idx = len(father.vals)
				}
				newvals := make([]Hasher, len(father.vals)+1)
				be := father.vals[:idx]
				ed := father.vals[idx:]
				copy(newvals[:idx], be)
				copy(newvals[idx+1:], ed)
				newvals[idx] = myint(ri.vals[0].Hash())
				father.vals = newvals
				childs := make([]*node, len(father.childs)+1)
				copy(childs[:idx], father.childs[:idx])
				childs[idx] = lf
				childs[idx+1] = ri
				copy(childs[idx+2:], father.childs[idx+1:])
				father.childs = childs
				if len(n.childs) != 0 { // 向上分裂
					ri.vals = ri.vals[1:]
				}
				n = father
				goto START
			} else if len(n.childs) != 0 { // 向上分裂
				ri.vals = ri.vals[1:]
				t.root = &node{
					vals:   []Hasher{myint(n.vals[t.m/2].Hash())},
					childs: []*node{lf, ri},
				}
				lf.father = t.root
				ri.father = t.root
			} else {
				n.vals = []Hasher{myint(ri.vals[0].Hash())}
				n.childs = append(n.childs, lf, ri)
			}
		}
		return
	}
	idx := -1
	for i, v := range n.vals {
		if v.Hash() >= val.Hash() {
			idx = i
			break
		}
	}
	if idx == -1 {
		idx = len(n.vals)
	}
	if len(n.childs) <= idx {
		n.childs = append(n.childs, &node{
			father: n,
		})
	}
	if len(n.childs) <= idx {
		println()
	}
	n.childs[idx].insert(t, val)
}
func (n *node) ensureReversePointer() {
	for _, v := range n.childs {
		v.father = n
	}
}

func (t *Tree) Travel(job func(val Hasher, level int)) {
	if t.root == nil {
		return
	}
	t.root.travel(job, 1)
}
func (n *node) travel(job func(val Hasher, level int), level int) {
	idx := 0
	for _, v := range n.childs {
		v.travel(job, level+1)
	}
	if len(n.childs) == 0 {
		for i, v := range n.vals {
			if i >= idx {
				job(v, level)
			}
		}
	}
}
