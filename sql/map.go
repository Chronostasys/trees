// map relational query to kv
package sql

import (
	"fmt"
	"strconv"
)

const (
	tablePrefix = "t_"
	tNumKey     = "tnum"
)

var (
	tree = makekv()
)

func CreateTable(table string) {
	num := GetTableNum()
	num++
	idstr := strconv.Itoa(num)
	tree.Insert(tNumKey, idstr)
	tree.Insert(fmt.Sprintf("%s%s", tablePrefix, table), "")
}
func DeleteTable(table string) {
	num := GetTableNum()
	num--
	idstr := strconv.Itoa(num)
	tree.Insert(tNumKey, idstr)
	tree.Delete(fmt.Sprintf("%s%s", tablePrefix, table))
}

// GetTableNum return table counts
func GetTableNum() int {
	idv := tree.Search(tNumKey)
	id, err := strconv.Atoi(idv)
	if err != nil {
		return 0
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
