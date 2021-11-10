package bitree

var (
	colorReset = "\033[0m"

	colorRed = "\033[31m"
)

type node struct {
	left   *node
	right  *node
	father *node
	val    Hasher
	red    bool
}

type Hasher interface {
	Hash() int
}

type Tree struct {
	root  *node
	total int
	level int
}

func (t *node) travel(job func(val Hasher)) {
	if t.left != nil {
		t.left.travel(job)
	}
	job(t.val)
	if t.right != nil {
		t.right.travel(job)
	}

}
func (t *node) print(matrix [][]*node, x, y, half int) {
	matrix[y][x] = t
	if t.left != nil {
		t.left.print(matrix, x-half, y+1, half/2)
	}
	if t.right != nil {
		t.right.print(matrix, x+half, y+1, half/2)
	}
}
func (t *node) search(hash int) Hasher {
	if hash == t.val.Hash() {
		return t.val
	}
	if hash < t.val.Hash() {
		if t.left == nil {
			return nil
		}
		return t.left.search(hash)
	} else {
		if t.right == nil {
			return nil
		}
		return t.right.search(hash)
	}
}
func (t *node) insert(val Hasher, tree *Tree, level int) {
	if val.Hash() < t.val.Hash() {
		if t.left == nil {
			t.left = &node{val: val, red: true, father: t}
			t.left.prebalance(tree).rebalance(tree)
			tree.level = level + 1
			return
		}
		t.left.insert(val, tree, level+1)
	} else {
		if t.right == nil {
			t.right = &node{val: val, red: true, father: t}
			t.right.prebalance(tree).rebalance(tree)
			tree.level = level + 1
			return
		}
		t.right.insert(val, tree, level+1)
	}
}

func (t *node) prebalance(tree *Tree) *node {
	if t.father.father != nil {
		if t.red && t.father.red && t.father.right == t && t.father.father.left == t.father {
			father := t.father
			father.rotateleft(tree)
			return father
		} else if t.red && t.father.red && t.father.left == t && t.father.father.right == t.father {
			father := t.father
			father.rotateright(tree)
			return father
		}
	}
	return t
}
func (t *node) setleft(n *node) {
	t.left = n
	if n == nil {
		return
	}
	n.father = t
}
func (t *node) setright(n *node) {
	t.right = n
	if n == nil {
		return
	}
	n.father = t
}
func (t *node) rotateleft(tree *Tree) {
	if t == tree.root {
		tree.root = t.right
	} else {
		if t.father.left == t {
			t.father.left = t.right
		} else {
			t.father.right = t.right
		}
	}
	t.right.father = t.father
	newRoot := t.right
	t.setright(newRoot.left)
	newRoot.setleft(t)
}
func (t *node) rotateright(tree *Tree) {
	if t == tree.root {
		tree.root = t.left
	} else {
		if t.father.left == t {
			t.father.left = t.left
		} else {
			t.father.right = t.left
		}
	}
	t.left.father = t.father
	newRoot := t.left
	t.setleft(newRoot.right)
	newRoot.setright(t)
}

func (t *node) rebalance(tree *Tree) {
	grandfa := t.father.father
	if t.father.red && grandfa != nil {
		other := grandfa.left
		lf := true
		if grandfa.left == t.father {
			other = grandfa.right
			lf = false
		}
		if other == nil || !other.red {
			grandfa.red = true
			t.father.red = false
			if lf {
				grandfa.rotateleft(tree)
			} else {
				grandfa.rotateright(tree)
			}
		} else {
			other.red = false
			t.father.red = false
			if tree.root != grandfa {
				grandfa.red = true
				if grandfa.father.red {
					grandfa.prebalance(tree).rebalance(tree)
				}
			}
		}
	}
}

func (t *Tree) Insert(val Hasher) {
	defer func() {
		t.total++
	}()
	if t.root == nil {
		t.root = &node{val: val}
		return
	}
	t.root.insert(val, t, 0)
}

func (t *Tree) Travel(job func(val Hasher)) {
	if t.root == nil {
		return
	}
	t.root.travel(job)
}
func (t *Tree) Print() {
	if t.root == nil {
		return
	}
	matrix := make([][]*node, t.level+1)
	for i := range matrix {
		matrix[i] = make([]*node, 1<<(t.level+2))
	}
	t.root.print(matrix, (1<<(t.level+2))/2, 0, (1<<(t.level+2))>>2)
	for _, v := range matrix {
		for _, u := range v {
			if u == nil {
				print(" ")
			} else {
				color := colorReset
				if u.red {
					color = colorRed
				}
				print(color, u.val.Hash())
			}
		}
		println(colorReset)
	}
}
func (t *Tree) Search(hash int) Hasher {
	if t.root == nil {
		return nil
	}
	return t.root.search(hash)
}
