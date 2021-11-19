package btree

import "testing"

func TestTree_Insert(t *testing.T) {
	tree := Make(3)
	tree.Insert(myint(1))
	tree.Insert(myint(2))
	tree.Insert(myint(3))
	tree.Insert(myint(4))
	tree.Insert(myint(5))
	tree.Insert(myint(6))
	println()
}
