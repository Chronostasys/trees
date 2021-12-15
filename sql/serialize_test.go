package sql

import (
	"reflect"
	"testing"
)

type Test struct {
	TestInt    int    `sql:"pk"`
	TestString string `idx:"-"`
	TestFloat  float32
}

func Test_serialize(t *testing.T) {
	Register(&Test{})
	s := &Test{TestInt: 9, TestString: "dafdsf", TestFloat: 1.1}
	bs := serialize(s, nil)
	test := &Test{}
	t.Run("Test deserialize", func(t *testing.T) {
		deserialize(bs, test)
		if !reflect.DeepEqual(test, s) {
			t.Errorf("value not equal after serialize and deserialize. before: %v after: %v", s, test)
		}
	})
	test = &Test{}
	t.Run("Test deserialize select", func(t *testing.T) {
		deserialize(bs, test, "TestString")
		exp := &Test{TestString: s.TestString}
		if !reflect.DeepEqual(test, exp) {
			t.Errorf("value not equal after serialize and selected deserialize. expect: %v got: %v", exp, test)
		}
	})
	t.Run("Test getpk", func(t *testing.T) {
		pk := metaMap[reflect.TypeOf(*s).String()].getpk(s)
		if pk != string(itb(int64(s.TestInt))) {
			t.Errorf("expect pk=%v, got %v", itb(int64(s.TestInt)), []byte(pk))
		}
	})
}
func BenchmarkSerialize(b *testing.B) {
	Register(&Test{})
	s := &Test{TestInt: 9, TestString: "dafdsf", TestFloat: 1.1}
	b.Run("BenchmarkSerialize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			serialize(s, nil)
		}
	})
	bs := serialize(s, nil)
	test := &Test{}
	b.Run("BenchmarkDeserialize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			deserialize(bs, test)
		}
	})
}
