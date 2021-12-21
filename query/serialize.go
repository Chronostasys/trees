package query

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

var (
	metaMap = map[string]SeriMeta{}
)

type SeriMeta struct {
	FieldsN  int
	Idx2Name map[int]fieldMeta
	Name2Idx map[string]int
	getpk    func(val interface{}) string
	Name     string
	Idx      map[string]int
	PKKind   reflect.Kind
	PKIdx    int
	t        reflect.Type
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

func getFieldStr(v reflect.Value, idx int) string {
	val := v.Field(idx)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return string(itb(val.Int()))
	case reflect.String:
		return val.String()
	case reflect.Float32, reflect.Float64:
		return string(ftb(val.Float()))
	}
	return ""
}

func (meta *SeriMeta) buildGetPK(i interface{}) {
	meta.getpk = func(i interface{}) string {
		v := reflect.Indirect(reflect.ValueOf(i))
		val := v.Field(meta.PKIdx)
		switch meta.PKKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return string(itb(val.Int()))
		case reflect.String:
			return val.String()
		case reflect.Float32, reflect.Float64:
			return string(ftb(val.Float()))
		}

		return ""
	}
}

func GetMeta(i interface{}) SeriMeta {
	v := reflect.Indirect(reflect.ValueOf(i))
	idxmap := map[int]fieldMeta{}
	tp := v.Type()
	pkidx := -1
	pkkind := reflect.Kind(0)
	idx := map[string]int{}
	name2Idx := map[string]int{}
	for i := 0; i < v.NumField(); i++ {
		idxmap[i] = fieldMeta{name: tp.Field(i).Name, kind: v.Field(i).Kind()}
		if tp.Field(i).Tag.Get("sql") == "pk" {
			pkidx = i
			pkkind = v.Field(i).Kind()
		}
		if len(tp.Field(i).Tag.Get("idx")) > 0 {
			idx[tp.Field(i).Name] = i
		}
		name2Idx[tp.Field(i).Name] = i
	}
	meta := SeriMeta{
		FieldsN:  v.NumField(),
		Idx2Name: idxmap,
		Name:     v.Type().String(),
		Idx:      idx,
		Name2Idx: name2Idx,
		PKKind:   pkkind,
		PKIdx:    pkidx,
		t:        reflect.TypeOf(i).Elem(),
	}
	meta.buildGetPK(i)
	return meta
}

func Register(i interface{}) {
	v := GetMeta(i)
	metaMap[v.Name] = v
}

func serialize(i interface{}, fmap map[int]func(s string)) []byte {
	v := reflect.Indirect(reflect.ValueOf(i))
	meta := metaMap[v.Type().String()]
	fieldsN := meta.FieldsN
	enc := []byte{}
	for i := 0; i < fieldsN; i++ {
		val := v.Field(i)
		var bs []byte
		switch meta.Idx2Name[i].kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bs = itb(val.Int())
		case reflect.String:
			s := val.String()
			enc = append(enc, itb(int64(len(s)))...)
			bs = []byte(s)
		case reflect.Float32, reflect.Float64:
			bs = ftb(val.Float())
		default:
			continue
		}
		enc = append(enc, bs...)
		if fmap != nil {
			if f, ok := fmap[i]; ok {
				f(string(bs))
			}

		}

	}
	return enc
}

var errDeserialize = fmt.Errorf("deserialize error")
var emptyMap = map[int]struct{}{}

func IterSerFields(ser []byte, meta SeriMeta, callback func(i int, v string)) {
	idx := 0
	le := len(ser)
	for i := 0; i < meta.FieldsN; i++ {
		fieldmeta := meta.Idx2Name[i]
		var bs []byte

		switch fieldmeta.kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bs = ser[idx : idx+8]
			idx += 8
		case reflect.String:
			if idx+8 > le {
				return
			}
			l := int(binary.LittleEndian.Uint64(ser[idx : idx+8]))
			idx += 8
			if idx+l > le {
				return
			}
			bs = ser[idx : idx+l]

			idx += l
		case reflect.Float32, reflect.Float64:
			if idx+8 > le {
				return
			}
			bs = ser[idx : idx+8]
			idx += 8
		}
		callback(i, string(bs))
	}
}

func deserializeEQ(ser []byte, i, into interface{}, eqfields map[int]struct{}, fields ...string) (succ bool, err error) {
	v := reflect.Indirect(reflect.ValueOf(i))
	vi := getIndirect(into)
	meta := metaMap[v.Type().String()]
	fieldsN := meta.FieldsN
	idx := 0
	m := make(map[string]struct{}, len(fields))
	for _, v := range fields {
		m[v] = struct{}{}
	}
	for i := 0; i < fieldsN; i++ {
		val := v.Field(i)
		fieldmeta := meta.Idx2Name[i]
		set := true
		if len(fields) != 0 {
			if _, ok := m[fieldmeta.name]; !ok {
				set = false
			}
		}
		vset := vi.Field(i)
		le := len(ser)
		switch fieldmeta.kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if set {
				if idx+8 > le {
					return false, errDeserialize
				}
				s := int64(binary.LittleEndian.Uint64(ser[idx : idx+8]))
				if _, ok := eqfields[i]; ok {
					if val.Int() != s {
						return false, nil
					}
				} else {
					vset.SetInt(s)
				}
			}

			idx += 8
		case reflect.String:
			if idx+8 > le {
				return false, errDeserialize
			}
			l := int(binary.LittleEndian.Uint64(ser[idx : idx+8]))
			idx += 8
			if set {
				if idx+l > le {
					return false, errDeserialize
				}
				s := string(ser[idx : idx+l])
				if _, ok := eqfields[i]; ok {
					if val.String() != s {
						return false, nil
					}
				} else {
					vset.SetString(s)
				}

			}

			idx += l
		case reflect.Float32, reflect.Float64:
			if set {
				if idx+8 > le {
					return false, errDeserialize
				}
				s := math.Float64frombits(binary.LittleEndian.Uint64(ser[idx : idx+8]))
				if _, ok := eqfields[i]; ok {
					if val.Float() != s {
						return false, nil
					}
				} else {
					vset.SetFloat(s)
				}
			}
			idx += 8
		}

	}
	return true, nil
}

func deserialize(ser []byte, i interface{}, fields ...string) error {
	_, err := deserializeEQ(ser, i, i, emptyMap, fields...)
	return err
}
