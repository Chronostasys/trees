// map relational query to kv
package sql

import (
	"reflect"
	"testing"
)

func TestTable(t *testing.T) {
	tableNames := []interface{}{"test_table1", 1}
	for _, v := range tableNames {
		CreateTable(v)
	}
	t.Run("TestGetTableMaxID", func(t *testing.T) {
		max := GetTableMaxID()
		if max != len(tableNames)-1 {
			t.Errorf("expect table maxid %d, got %d", len(tableNames)-1, max)
		}
	})
	t.Run("TestGetTableNames", func(t *testing.T) {
		names := GetTableNames()
		exp := []string{"int", "string"}
		if !reflect.DeepEqual(names, exp) {
			t.Errorf("expect tablenames %v, got %v", exp, names)
		}
	})
	t.Run("TestDeleteTable", func(t *testing.T) {
		DeleteTable(tableNames[0])
		names := GetTableNames()
		if !reflect.DeepEqual([]string{"int"}, names) {
			t.Errorf("expect tablenames %v, got %v", []string{"int"}, names)
		}
	})
}

func TestRow(t *testing.T) {
	Register(&Test{})
	CreateTable(&Test{})
	q, _ := Table(&Test{})
	item := &Test{
		TestInt:    11,
		TestString: "test",
		TestFloat:  9.33,
	}
	q.Insert(&Test{
		TestInt:    10,
		TestString: "atest",
		TestFloat:  9.33,
	})
	item2 := &Test{
		TestInt:    12,
		TestString: "test",
		TestFloat:  9.34,
	}
	q.Insert(item2)
	q.Insert(item)
	t.Run("test find by pk", func(t *testing.T) {
		re := &Test{TestInt: 11}
		q.FindByPK(re)
		if !reflect.DeepEqual(re, item) {
			t.Errorf("expect search result %v, got %v", item, re)
		}
	})
	t.Run("test find by pk not found", func(t *testing.T) {
		re := &Test{TestInt: 9}
		err := q.FindByPK(re)
		if err != errNotfound {
			t.Errorf("expect err not found, got %v", err)
		}
	})
	t.Run("test find by index", func(t *testing.T) {
		re := &Test{TestString: "test"}
		err := q.Find(re, "TestString")
		if !reflect.DeepEqual(re, item) {
			t.Errorf("expect search result %v, got %v. err=%v", item, re, err)
		}
	})
	t.Run("test find by multi query with single index", func(t *testing.T) {
		re := &Test{TestString: "test", TestFloat: 9.34}
		err := q.Find(re, "TestString", "TestFloat")
		if !reflect.DeepEqual(re, item2) {
			t.Errorf("expect search result %v, got %v. err=%v", item2, re, err)
		}
	})
	t.Run("test find no index", func(t *testing.T) {
		re := &Test{TestFloat: 9.34}
		err := q.Find(re, "TestFloat")
		if !reflect.DeepEqual(re, item2) {
			t.Errorf("expect search result %v, got %v. err=%v", item2, re, err)
		}
	})
}
