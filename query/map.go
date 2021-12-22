// map relational query to kv
package query

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	tablePrefix       = "t_"
	rowTemplatePrefix = "i_%3d_r"   //主键一定不是空，没必要把_r去了
	rowTemplate       = "i_%3d_r%s" //i_{tableid}_r{rowpk}
	tNumKey           = "tnum"
	idxTemplatePrefix = "i_%3d_i%s_v-%s"
	idxTemplate       = "i_%3d_i%s_v-%s_p%s" //i_{tableid}_i{idx}_v{val}_p{pk}
	uniqueIdxTemplate = "i_%3d_i%s_v-%s"     //i_{tableid}_i{idx}_v{val}
	idxNumKey         = "inum"
)

var (
	tree           KV
	errNoSuchTable = fmt.Errorf("no such table")
	errNotfound    = fmt.Errorf("not found")
)

func SetKV(kv KV) {
	tree = kv
}

func getIndirect(i interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(i))
}
func getTypeName(t interface{}) string {
	return getIndirect(t).Type().String()
}
func CreateTable(t interface{}) int {
	key := fmt.Sprintf("%s%s", tablePrefix, getTypeName(t))
	ext := tree.Search(key)
	if i, err := strconv.Atoi(ext); err == nil { // exist
		return i
	}
	id := GetTableMaxID()
	id++
	idstr := strconv.Itoa(id)
	tree.Insert(tNumKey, idstr)
	tree.Insert(key, idstr)
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
	tree.Larger(tablePrefix, 2000, 1000, 0, func(k, v string) bool {
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
	meta SeriMeta
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
	fmap := make(map[int]func(s string), len(meta.Idx))
	for _, v := range meta.Idx {
		fmap[v] = func(s string) {
			tree.Insert(fmt.Sprintf(idxTemplate, q.tid, string(itb(int64(v))), s, k), "")
		}
	}
	tree.Insert(k, string(serialize(i, fmap)))
}
func (q *TableQuerier) Update(i interface{}, fields ...string) {
	meta := q.meta
	k := fmt.Sprintf(rowTemplate, q.tid, meta.getpk(i))
	fmap := make(map[int]func(s string), len(meta.Idx))
	oldidx := map[int]struct{}{}
	for _, v := range fields {
		if i, ok := meta.Idx[v]; ok {
			oldidx[i] = struct{}{}
		}
	}
	IterSerFields([]byte(tree.Search(k)), meta, func(i int, v string) {
		if _, ok := oldidx[i]; ok {
			k1 := fmt.Sprintf(idxTemplate, q.tid, string(itb(int64(i))), v, k)
			tree.Delete(k1)
		}
	})
	for _, v := range meta.Idx {
		fmap[v] = func(s string) {
			tree.Insert(fmt.Sprintf(idxTemplate, q.tid, string(itb(int64(v))), s, k), "")
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
func (q *TableQuerier) FindOne(i interface{}, fields ...string) error {
	return q.Find(i, 0, 1, func(i interface{}) bool {
		return false
	}, fields...)
}
func appendToSlice(arrPtr, val interface{}) {
	valuePtr := reflect.ValueOf(arrPtr)
	value := valuePtr.Elem()

	value.Set(reflect.Append(value, reflect.ValueOf(val)))
}

func (q *TableQuerier) FindAll(i interface{}, arrptr interface{}, skip, limit int, fields ...string) error {
	return q.Find(i, 0, limit, func(i interface{}) bool {
		appendToSlice(arrptr, i)
		return true
	}, fields...)
}
func (q *TableQuerier) Find(i interface{}, skip, limit int, callback func(i interface{}) bool, fields ...string) error {
	meta := q.meta
	idx := -1
	m := make(map[int]struct{}, len(fields)-1)
	for _, v := range fields {
		if index, ok := meta.Idx[v]; ok && idx == -1 {
			idx = index
		} else {
			m[meta.Name2Idx[v]] = struct{}{}
		}
	}
	v := getIndirect(i)
	searchLimit := limit
	if len(fields) != idx {
		searchLimit = 1000
	}

	if idx != -1 { // use index
		idxprefix := fmt.Sprintf(idxTemplatePrefix, q.tid, string(itb(int64(idx))), getFieldStr(v, idx))
		succ := false
		tree.Larger(idxprefix, 2000, searchLimit, skip, func(k, v string) bool {
			if len(k) <= len(idxprefix) || k[:len(idxprefix)] != idxprefix {
				return false
			}
			ser := tree.Search(k[len(idxprefix)+2:])
			into := i
			if limit != 1 {
				into = reflect.New(meta.t).Interface()
			}
			succ, _ = deserializeEQ([]byte(ser), i, into, m)
			if succ {
				return callback(into)
			}
			return !succ
		})
		if !succ {
			return errNotfound
		}
		return nil
	}
	idxprefix := fmt.Sprintf(rowTemplatePrefix, q.tid)
	succ := false
	tree.Larger(idxprefix, 2000, searchLimit, skip, func(k, v string) bool {
		if len(k) <= len(idxprefix) || k[:len(idxprefix)] != idxprefix {
			return false
		}
		into := i
		if limit != 1 {
			into = reflect.New(meta.t).Interface()
		}
		succ, _ = deserializeEQ([]byte(v), i, into, m)
		if succ {
			return callback(into)
		}
		return !succ
	})
	if !succ {
		return errNotfound
	}
	return nil
}
