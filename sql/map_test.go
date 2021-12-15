// map relational query to kv
package sql

import (
	"reflect"
	"testing"
)

func TestTable(t *testing.T) {
	tableNames := []string{"test_table1", "test_table2"}
	for _, v := range tableNames {
		CreateTable(v)
	}
	t.Run("TestGetTableNum", func(t *testing.T) {
		max := GetTableMaxID()
		if max != len(tableNames)-1 {
			t.Errorf("expect table maxid %d, got %d", len(tableNames)-1, max)
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

func TestRow(t *testing.T) {
	Register(&Test{})
	tableNames := []string{"test_table1", "test_table2"}
	for _, v := range tableNames {
		CreateTable(v)
	}
	q, _ := Table(tableNames[0])
	item := &Test{
		TestInt:    11,
		TestString: "test",
		TestFloat:  9.33,
	}
	q.Insert(item)
	t.Run("test find by pk", func(t *testing.T) {
		re := &Test{TestInt: 11}
		q.FindByPK(re)
		if !reflect.DeepEqual(re, item) {
			t.Errorf("expect search result %v, got %v", item, re)
		}
	})
	t.Run("test find by pk not found", func(t *testing.T) {
		re := &Test{TestInt: 10}
		err := q.FindByPK(re)
		if err != errNotfound {
			t.Errorf("expect err not found, got %v", err)
		}
	})
}
