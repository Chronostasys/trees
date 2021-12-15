// map relational query to kv
package sql

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	tablePrefix = "t_"
	rowTemplate = "r_%d_%s" //r_tableid_rowpk
	tNumKey     = "tnum"
)

var (
	tree           = makekv()
	errNoSuchTable = fmt.Errorf("no such table")
	errNotfound    = fmt.Errorf("not found")
)

func CreateTable(table string) int {
	id := GetTableMaxID()
	id++
	idstr := strconv.Itoa(id)
	tree.Insert(tNumKey, idstr)
	tree.Insert(fmt.Sprintf("%s%s", tablePrefix, table), idstr)
	return id
}
func DeleteTable(table string) {
	tree.Delete(fmt.Sprintf("%s%s", tablePrefix, table))
}
func GetTableID(table string) (id int, err error) {
	id, err = strconv.Atoi(tree.Search(fmt.Sprintf("%s%s", tablePrefix, table)))
	if err != nil {
		err = errNoSuchTable
	}
	return
}

// GetTableMaxID return table counts
func GetTableMaxID() int {
	idv := tree.Search(tNumKey)
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
	tid int
}

func Table(table string) (*TableQuerier, error) {
	id, err := GetTableID(table)
	if err != nil {
		return nil, err
	}
	return &TableQuerier{
		tid: id,
	}, nil
}

func (q *TableQuerier) Insert(i interface{}) {
	v := reflect.Indirect(reflect.ValueOf(i))
	meta := metaMap[v.Type().String()]
	k := fmt.Sprintf(rowTemplate, q.tid, meta.getpk(v))
	tree.Insert(k, string(serialize(i)))
}
func (q *TableQuerier) FindByPK(i interface{}, fields ...string) error {
	v := reflect.Indirect(reflect.ValueOf(i))
	meta := metaMap[v.Type().String()]
	k := fmt.Sprintf(rowTemplate, q.tid, meta.getpk(v))
	err := deserialize([]byte(tree.Search(k)), i, fields...)
	if err != nil {
		return errNotfound
	}
	return nil
}
