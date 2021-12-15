// map relational query to kv
package sql

import (
	"reflect"
	"testing"
)

func TestGetTableNames(t *testing.T) {
	tableNames := []string{"test_table1", "test_table2"}
	for _, v := range tableNames {
		CreateTable(v)
	}
	t.Run("TestGetTableNum", func(t *testing.T) {
		max := GetTableNum()
		if max != len(tableNames) {
			t.Errorf("expect table num %d, got %d", len(tableNames), max)
		}
	})
	t.Run("TestGetTableNames", func(t *testing.T) {
		names := GetTableNames()
		if !reflect.DeepEqual(tableNames, names) {
			t.Errorf("expect tablenames %v, got %v", tableNames, names)
		}
	})
	t.Run("TestDeleteTable", func(t *testing.T) {
		DeleteTable(tableNames[0])
		names := GetTableNames()
		if !reflect.DeepEqual(tableNames[1:], names) {
			t.Errorf("expect tablenames %v, got %v", tableNames[1:], names)
		}
	})
}
