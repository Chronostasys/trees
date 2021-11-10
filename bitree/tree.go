package bitree

import (
	"fmt"
)

type node struct {
	left  *node
	right *node
	// father *node
	val Hasher
	red bool
}

type Hasher interface {
	Hash() int
}

type Tree struct {
	root  *node
	total int
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
func (t *node) print(matrix [][]int, x, y, half int) {
	if y >= 10 {
		fmt.Println(y)
	}
	matrix[y][x] = t.val.Hash()
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
	defer func() {
		t.total++
	}()
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
func (t *Tree) Print() {
	if t.root == nil {
		return
	}
	matrix := make([][]int, t.total)
	for i := range matrix {
		matrix[i] = make([]int, 1<<t.total)
	}
	t.root.print(matrix, (1<<t.total)/2, 0, (1<<t.total)>>2)
	for _, v := range matrix {
		for _, u := range v {
			if u == 0 {
				print(" ")
			} else {
				print(u)
			}
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
