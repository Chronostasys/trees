package sql

import "github.com/Chronostasys/trees/btree"

type kvbtree struct {
	*btree.Tree
}

func makekv() *kvbtree {
	return &kvbtree{
		Tree: btree.Make(128),
	}
}

func (t *kvbtree) Insert(k, v string) {
	t.Tree.Insert(btree.KV{K: k, V: v})
}
func (t *kvbtree) Delete(k string) {
	t.Tree.Delete(btree.KV{K: k})
}
func (t *kvbtree) Search(k string) string {
	re := t.Tree.Search(btree.KV{K: k})
	if re == nil {
		return ""
	}
	return re.(btree.KV).V
}
func (t *kvbtree) Larger(than string, max int, callback func(k, v string) bool) {
	t.Tree.Larger(btree.KV{K: than}, max, func(i btree.Item) bool {
		kv := i.(btree.KV)
		return callback(kv.K, kv.V)
	})
}
func (t *kvbtree) LargerOrEq(than string, max int, callback func(k, v string) bool) {
	t.Tree.LargerOrEq(btree.KV{K: than}, max, func(i btree.Item) bool {
		kv := i.(btree.KV)
		return callback(kv.K, kv.V)
	})
}
