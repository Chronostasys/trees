package sql

import (
	"encoding/binary"
	"math"
	"reflect"
)

var (
	metaMap = map[reflect.Type]seriMeta{}
)

type seriMeta struct {
	fieldsN  int
	idx2Name map[int]fieldMeta
}
type fieldMeta struct {
	name string
	kind reflect.Kind
}

func itb(i int64) []byte {
	bs := [8]byte{}
	binary.LittleEndian.PutUint64(bs[:], uint64(i))
	return bs[:]
}
func ftb(i float64) []byte {
	bs := [8]byte{}
	binary.LittleEndian.PutUint64(bs[:], math.Float64bits(i))
	return bs[:]
}

func Register(i interface{}) {
	v := reflect.Indirect(reflect.ValueOf(i))
	idxmap := map[int]fieldMeta{}
	tp := v.Type()
	for i := 0; i < v.NumField(); i++ {
		idxmap[i] = fieldMeta{name: tp.Field(i).Name, kind: v.Field(i).Kind()}
	}
	meta := seriMeta{
		fieldsN:  v.NumField(),
		idx2Name: idxmap,
	}
	metaMap[reflect.TypeOf(i).Elem()] = meta
}

func serialize(i interface{}) []byte {
	v := reflect.Indirect(reflect.ValueOf(i))
	meta := metaMap[v.Type()]
	fieldsN := meta.fieldsN
	enc := []byte{}
	for i := 0; i < fieldsN; i++ {
		val := v.Field(i)
		switch meta.idx2Name[i].kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			enc = append(enc, itb(val.Int())...)
		case reflect.String:
			s := val.String()
			enc = append(enc, itb(int64(len(s)))...)
			enc = append(enc, []byte(s)...)
		case reflect.Float32, reflect.Float64:
			enc = append(enc, ftb(val.Float())...)
		}

	}
	return enc
}

func deserialize(ser []byte, i interface{}, fields ...string) {
	v := reflect.Indirect(reflect.ValueOf(i))
	meta := metaMap[v.Type()]
	fieldsN := meta.fieldsN
	idx := 0
	m := make(map[string]struct{}, len(fields))
	for _, v := range fields {
		m[v] = struct{}{}
	}
	for i := 0; i < fieldsN; i++ {
		val := v.Field(i)
		fieldmeta := meta.idx2Name[i]
		set := true
		if len(fields) != 0 {
			if _, ok := m[fieldmeta.name]; !ok {
				set = false
			}
		}
		switch fieldmeta.kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if set {
				val.SetInt(int64(binary.LittleEndian.Uint64(ser[idx : idx+8])))
			}

			idx += 8
		case reflect.String:
			l := int(binary.LittleEndian.Uint64(ser[idx : idx+8]))
			idx += 8
			if set {
				val.SetString(string(ser[idx : idx+l]))
			}

			idx += l
		case reflect.Float32, reflect.Float64:
			if set {
				val.SetFloat(math.Float64frombits(binary.LittleEndian.Uint64(ser[idx : idx+8])))
			}
			idx += 8
		}

	}
}
