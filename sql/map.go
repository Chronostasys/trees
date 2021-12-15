// map relational query to kv
package sql

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	tablePrefix       = "t_"
	rowTemplatePrefix = "i_%3d_r"
	rowTemplate       = "i_%3d_r%s" //i_{tableid}_r{rowpk}
	tNumKey           = "tnum"
	idxTemplate       = "i_%3d_i%s_v-%s_p%s" //i_{tableid}_i{idx}_v{val}_p{pk}
	uniqueIdxTemplate = "i_%3d_i%s_v-%s"     //i_{tableid}_i{idx}_v{val}
	idxNumKey         = "inum"
)

var (
	tree           = makekv()
	errNoSuchTable = fmt.Errorf("no such table")
	errNotfound    = fmt.Errorf("not found")
)

func getIndirect(i interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(i))
}
func getTypeName(t interface{}) string {
	return getIndirect(t).Type().String()
}
func CreateTable(t interface{}) int {
	id := GetTableMaxID()
	id++
	idstr := strconv.Itoa(id)
	tree.Insert(tNumKey, idstr)
	tree.Insert(fmt.Sprintf("%s%s", tablePrefix, getTypeName(t)), idstr)
	return id
}
func DeleteTable(t interface{}) {
	tree.Delete(fmt.Sprintf("%s%s", tablePrefix, getTypeName(t)))
}
func GetTableID(t interface{}) (id int, err error) {
	id, err = strconv.Atoi(tree.Search(fmt.Sprintf("%s%s", tablePrefix, getTypeName(t))))
	if err != nil {
		err = errNoSuchTable
	}
	return
}

// GetTableMaxID return table counts
func GetTableMaxID() int {
	return getNum(tNumKey)
}

func getNum(k string) int {
	idv := tree.Search(k)
	id, err := strconv.Atoi(idv)
	if err != nil {
		return -1
	}
	return id
}

func GetTableNames() []string {
	names := make([]string, 0)
	tree.Larger(tablePrefix, 1000, func(k, v string) bool {
		if len(k) <= len(tablePrefix) || k[:len(tablePrefix)] != tablePrefix {
			return false
		}
		names = append(names, k[len(tablePrefix):])
		return true
	})
	return names
}

type TableQuerier struct {
	tid  int
	meta seriMeta
}

func Table(t interface{}) (*TableQuerier, error) {
	id, err := GetTableID(t)
	if err != nil {
		return nil, err
	}
	return &TableQuerier{
		tid:  id,
		meta: metaMap[reflect.Indirect(reflect.ValueOf(t)).Type().String()],
	}, nil
}

func (q *TableQuerier) Insert(i interface{}) {
	meta := q.meta
	k := fmt.Sprintf(rowTemplate, q.tid, meta.getpk(i))
	fmap := map[int]func(s string){}
	for _, v := range meta.idx {
		fmap[v] = func(s string) {
			tree.Insert(fmt.Sprintf(idxTemplate, q.tid, string(itb(int64(v))), s, k), k)
		}
	}
	tree.Insert(k, string(serialize(i, fmap)))
}
func (q *TableQuerier) FindByPK(i interface{}, selfields ...string) error {
	meta := q.meta
	k := fmt.Sprintf(rowTemplate, q.tid, meta.getpk(i))
	err := deserialize([]byte(tree.Search(k)), i, selfields...)
	if err != nil {
		return errNotfound
	}
	return nil
}
func (q *TableQuerier) Find(i interface{}, fields ...string) error {
	meta := q.meta
	idx := -1
	m := map[int]struct{}{}
	for _, v := range fields {
		if index, ok := meta.idx[v]; ok && idx == -1 {
			idx = index
		} else {
			m[meta.name2Idx[v]] = struct{}{}
		}
	}
	v := getIndirect(i)
	if idx != -1 { // use index
		idxprefix := fmt.Sprintf(uniqueIdxTemplate, q.tid, string(itb(int64(idx))), getFieldStr(v, idx))
		succ := false
		tree.Larger(idxprefix, 1000, func(k, v string) bool {
			if len(k) <= len(idxprefix) || k[:len(idxprefix)] != idxprefix {
				return false
			}
			ser := tree.Search(v)
			succ, _ = deserializeEQ([]byte(ser), i, m)
			return !succ
		})
		if !succ {
			return errNotfound
		}
		return nil
	}
	idxprefix := fmt.Sprintf(rowTemplatePrefix, q.tid)
	succ := false
	tree.Larger(idxprefix, 1000, func(k, v string) bool {
		if len(k) <= len(idxprefix) || k[:len(idxprefix)] != idxprefix {
			return false
		}
		succ, _ = deserializeEQ([]byte(v), i, m)
		return !succ
	})
	if !succ {
		return errNotfound
	}
	return nil
}
