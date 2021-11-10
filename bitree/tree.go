package bitree

type node struct {
	left  *node
	right *node
	// father *node
	val Hasher
}

type Hasher interface {
	Hash() int
}

type Tree struct {
	root *node
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
func (t *node) insert(val Hasher) {
	if val.Hash() < t.val.Hash() {
		if t.left == nil {
			t.left = &node{val: val}
			return
		}
		t.left.insert(val)
	} else {
		if t.right == nil {
			t.right = &node{val: val}
			return
		}
		t.right.insert(val)
	}
}

func (t *Tree) Insert(val Hasher) {
	if t.root == nil {
		t.root = &node{val: val}
		return
	}
	t.root.insert(val)
}

func (t *Tree) Travel(job func(val Hasher)) {
	if t.root == nil {
		return
	}
	t.root.travel(job)
}
func (t *Tree) Search(hash int) Hasher {
	if t.root == nil {
		return nil
	}
	return t.root.search(hash)
}
