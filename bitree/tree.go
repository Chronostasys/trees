package bitree

import "fmt"

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
}

func (t *node) travel(job func(val Hasher, level int), level int) {
	if t.left != nil {
		t.left.travel(job, level+1)
	}
	job(t.val, level)
	if t.right != nil {
		t.right.travel(job, level+1)
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

func (t *node) deleteMe(tree *Tree) {
	var n *node
	if t.left == nil && t.right == nil {
		n = t
		if n.red {
			if n.father.left == n {
				n.father.left = nil
			} else {
				n.father.right = nil
			}
			return
		}
	} else if t.left == nil && t.right != nil {
		n = t.right
		if n.red {
			t.val = n.val
			t.left = n.left
			t.right = n.right
			if t.left != nil {
				t.left.father = t
			}
			if t.right != nil {
				t.right.father = t
			}
			return
		}
	} else if t.right == nil && t.left != nil {
		n = t.left
		if n.red {
			t.val = n.val
			t.left = n.left
			t.right = n.right
			if t.left != nil {
				t.left.father = t
			}
			if t.right != nil {
				t.right.father = t
			}
			return
		}
	} else {
		n = t.right
		for {
			if n.left == nil {
				break
			}
			n = n.left
		}
		if n.red {
			if n.father == t {
				t.right = n.right
				if n.right != nil {
					n.right.father = t
				}
			} else {
				n.father.left = n.right
				if n.right != nil {
					n.right.father = n.father
				}
			}
			t.val = n.val
			return
		}
	}
	n.deleteRedBlack(tree, t)

	t.val = n.val
	if n.father == nil {
		tree.root = nil
	} else {
		if n.father.left == n {
			n.father.setleft(n.right)
		} else {
			n.father.setright(n.right)
		}
	}
}
func (t *node) deleteRedBlack(tree *Tree, del *node) {
	if t.red {
		t.deleteRed(tree, del)
	} else {
		t.deleteBlack(tree, del)
	}
}
func (t *node) deleteRed(tree *Tree, del *node) {
	t.red = false
}

func (t *node) deleteBlack(tree *Tree, del *node) {
	if t.father == nil {
		return
	}
	if t.father.left == t && t.father.right != nil && t.father.right.red {
		t.father.right.red = false
		t.father.red = true
		t.father.rotateleft(tree)
		t.deleteRedBlack(tree, del)
		return
	}
	if t.father.right == t && t.father.left != nil && t.father.left.red {
		t.father.left.red = false
		t.father.red = true
		t.father.rotateright(tree)
		t.deleteRedBlack(tree, del)
		return
	}
	if t.father.left == t &&
		(t.father.right != nil && !t.father.right.red) {
		bro := t.father.right
		if bro.right != nil && bro.right.red {
			bro.red = bro.father.red
			bro.right.red = false
			bro.father.red = false
			bro.father.rotateleft(tree)
			// t.deleteRedBlack(tree, del)
			return
		}
		if bro.left != nil && bro.left.red {
			bro.left.red = false
			bro.red = true
			bro.rotateright(tree)
			t.deleteRedBlack(tree, del)
			return
		}
		bro.red = true
		bro.father.deleteRedBlack(tree, del)
		return
	}
	if t.father.right == t &&
		(t.father.left != nil && !t.father.left.red) {
		bro := t.father.left

		if bro.left != nil && bro.left.red {
			bro.red = bro.father.red
			bro.left.red = false
			bro.father.red = false
			bro.father.rotateright(tree)
			// t.deleteRedBlack(tree, del)
			return
		}
		if bro.right != nil && bro.right.red {
			bro.right.red = false
			bro.red = true
			bro.rotateleft(tree)
			t.deleteRedBlack(tree, del)
			return
		}
		bro.red = true
		bro.father.deleteRedBlack(tree, del)
		return
	}
}
func (t *node) delete(hash int, tree *Tree) {
	if hash == t.val.Hash() {
		t.deleteMe(tree)
		tree.total--
		return
	}
	if hash < t.val.Hash() {
		if t.left == nil {
			return
		}
		t.left.delete(hash, tree)
		return
	} else {
		if t.right == nil {
			return
		}
		t.right.delete(hash, tree)
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
func (t *node) insert(val Hasher, tree *Tree, level int) bool {
	if val.Hash() < t.val.Hash() {
		if t.left == nil {
			t.left = &node{val: val, red: true, father: t}
			t.left.prebalance(tree).rebalance(tree)
			return true
		}
		t.left.insert(val, tree, level+1)
		return true
	} else if val.Hash() > t.val.Hash() {
		if t.right == nil {
			t.right = &node{val: val, red: true, father: t}
			t.right.prebalance(tree).rebalance(tree)
			return true
		}
		t.right.insert(val, tree, level+1)
		return true
	} else {
		t.val = val
		return false
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
	if t.root == nil {
		t.root = &node{val: val}
		return
	}
	newv := t.root.insert(val, t, 0)
	if newv {
		t.total++
	}
}

func (t *Tree) Travel(job func(val Hasher, level int)) {
	if t.root == nil {
		return
	}
	t.root.travel(job, 1)
}
func (t *Tree) Print(colored bool) {
	if t.root == nil {
		return
	}
	level := 0
	t.Travel(func(val Hasher, lev int) {
		if level < lev {
			level = lev
		}
	})
	matrix := make([][]*node, level)
	for i := range matrix {
		matrix[i] = make([]*node, 1<<(level))
	}
	t.root.print(matrix, (1<<(level))/2, 0, (1<<(level))>>2)
	for _, v := range matrix {
		for _, u := range v {
			if u == nil {
				print("  ")
			} else {
				if colored {
					color := colorReset
					if u.red {
						color = colorRed
					}
					print(color)
				}
				print(fmt.Sprintf("%02d", u.val.Hash()))
			}
		}
		if colored {
			print(colorReset)
		}
		println()
	}
}
func (t *Tree) Search(hash int) Hasher {
	if t.root == nil {
		return nil
	}
	return t.root.search(hash)
}

func (t *Tree) Delete(hash int) {
	if t.root == nil {
		return
	}
	t.root.delete(hash, t)
}

func (t *Tree) Len() int {
	return t.total
}

func Make() *Tree {
	return &Tree{}
}
