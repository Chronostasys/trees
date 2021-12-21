package query

import "github.com/Chronostasys/trees/btree"

type kvbtree struct {
	*btree.Tree
}

type KV interface {
	Insert(k, v string)
	Delete(k string)
	Search(k string) string
	Larger(than string, max, limit, skip int, callback func(k, v string) bool)
}

func MakeLocalInMemKV() KV {
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
func (t *kvbtree) Larger(than string, max, limit, skip int, callback func(k, v string) bool) {
	t.Tree.Larger(btree.KV{K: than}, max, limit, skip, func(i btree.Item) bool {
		kv := i.(btree.KV)
		return callback(kv.K, kv.V)
	})
}
