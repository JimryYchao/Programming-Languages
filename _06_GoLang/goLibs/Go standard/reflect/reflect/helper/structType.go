package helper

import (
	"fmt"
	r "reflect"
	"strconv"
)

//! >>>>>>>>>>>> Fields <<<<<<<<<<<<

type Fields []Field

func (fs Fields) ToStructFields() []r.StructField {
	fields := make([]r.StructField, len(fs))
	for i, v := range fs {
		fields[i] = v.StructField
	}
	return fields
}
func (fs Fields) Add(fields ...Field) Fields {
	newfs := append(fs, fields...)
	return newfs
}
func (fs Fields) AddStructField(fields ...r.StructField) Fields {
	newfs := make(Fields, len(fields))
	for i, v := range fields {
		newfs[i] = &field{v, nil}
	}
	return fs.Add(newfs...)
}

//! >>>>>>>>>>>> Field <<<<<<<<<<<<

type Field = *field

type field struct {
	r.StructField
	keyValues map[string]string
}

func (f Field) Type() Type {
	return TypeFrom(f.StructField.Type)
}

func (f Field) Get(key string) string {
	if f.keyValues == nil {
		f.TagKeys()
	}
	return f.keyValues[key]
}
func (f Field) TagKeys() []string {
	tag := string(f.Tag)
	count := 0
	keys := make([]string, 0, 10)
	if f.keyValues == nil {
		f.keyValues = make(map[string]string)
	}
	clear(f.keyValues)
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			return nil
		}
		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			return nil
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			return nil
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			return nil
		}

		keys = append(keys, name)
		f.keyValues[name] = value
		count++
	}
	if count > 0 {
		return keys[0:count]
	}
	clear(f.keyValues)
	return nil
}

func (f Field) String() string {
	return fmt.Sprintf("%v", f.StructField)
}

//! >>>>>>>>>>>> StructType <<<<<<<<<<<<

type StructType = *structType
type structType struct {
	*tCommon
	num    int
	fields []Field
}

func (t StructType) typeof(tp r.Type) Type {
	st := structType{}
	st.tCommon = newTCommon(tp)
	st.num = tp.NumField()
	st.fields = make([]Field, st.num)
	for i := range st.num {
		st.fields[i] = &field{tp.Field(i), nil}
	}
	t = &st
	return t
}

func (StructType) Kind() r.Kind    { return r.Struct }
func (t StructType) NumField() int { return t.num }

func (t StructType) Field(i int) (Field, bool) {
	if i < 0 || i >= t.num {
		return nil, false
	}
	return t.fields[i], true
}

func (t StructType) Fields() []Field { return t.fields }

func (t StructType) FieldByName(name string) (Field, bool) {
	if f, ok := t.t.FieldByName(name); ok {
		return &field{f, nil}, true
	}
	return nil, false
}

func (t StructType) FieldByNameFunc(match func(s string) bool) (Field, bool) {
	if f, ok := t.t.FieldByNameFunc(match); ok {
		return &field{f, nil}, true
	}
	return nil, false
}

func (t StructType) FieldByIndex(index []int) (Field, bool) {
	if len(index) == 0 {
		return nil, false
	}
	field, ok := t.Field(index[0])
	var f Type = field.Type()
	for _, x := range index[1:] {
		ft := f
		if p, ok := TypeTo(ft).PointerType(); ok && p.Elem().Kind() == r.Struct {
			ft = p.Elem()
		}
		f = ft
		if f.Kind() == r.Struct {
			st := f.(StructType)
			if field, ok = st.Field(x); !ok {
				return nil, false
			}
			f = field.Type()
		} else {
			return nil, false
		}
	}
	return field, ok
}

// StructOf

func StructOf(fields []r.StructField) (st StructType, e error) {
	defer func() {
		if err := recover(); err != nil {
			st = nil
			e = newErr("StructOf", err.(error))
		}
	}()
	st = new(structType)
	st = st.typeof(r.StructOf(fields)).(StructType)
	return
}

func StructOfFields(fields Fields) (StructType, error) {
	fs := fields.ToStructFields()
	return StructOf(fs)
}

func VisibleFields(t StructType) Fields {
	if t == nil {
		return nil
	}
	fs := r.VisibleFields(t.Type())
	fields := make([]Field, len(fs))
	for i, v := range fs {
		fields[i] = &field{v, nil}
	}
	return fields
}
