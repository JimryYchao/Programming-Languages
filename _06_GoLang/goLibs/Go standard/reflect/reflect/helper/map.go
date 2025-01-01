package helper

import (
	"fmt"
	r "reflect"
)

//! >>>>>>>>>>>>>> Map <<<<<<<<<<<<<<

type Map = *vMap

type vMap struct {
	*vCommon
}

func (v Map) valueof(rv r.Value) Value {
	v = &vMap{newValue(rv)}
	return v
}
func (v Map) Kind() r.Kind { return r.Map }
func (v Map) MapType() MapType {
	t, _ := TypeTo(v.Type()).MapType()
	return t
}

func (v Map) Len() int { return v.v.Len() }
func (v Map) Clear()   { v.v.Clear() }

func (v Map) MapKeys() []MKey {
	keys := v.v.MapKeys()
	outs := make([]MKey, len(keys))
	for i := range len(outs) {
		outs[i] = &mKey{ValueOf(keys[i])}
	}
	return outs
}
func (v Map) MapIndex(key Value) (Value, bool) {
	rt := v.v.MapIndex(key.Value())
	if rt.Kind() == r.Invalid {
		return nil, false
	}
	return newValue(rt), true
}
func (v Map) MapRange() *MapIter {
	return &MapIter{v.v.MapRange(), false, false}
}

func (v Map) SetMapIndex(key Value, value Value) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s", e)
		}
	}()
	var in = r.Value{}
	if value != nil && value.Kind() != r.Invalid {
		in = value.Value()
	}
	v.v.SetMapIndex(key.Value(), in)
	return nil
}
func (v Map) DeleteKey(key Value) bool {
	return v.SetMapIndex(key, nil) == nil
}

type MapIter struct {
	iter        *r.MapIter
	isStart     bool
	isExhausted bool
}
type (
	MKey   = *mKey
	MValue = *mValue
)
type mKey struct {
	Value
}
type mValue struct {
	Value
}

func (l *MapIter) Next() bool {
	if !l.isExhausted && l.iter.Next() {
		l.isStart = true
		return true
	}
	l.isExhausted = true
	return false
}
func (l *MapIter) Key() MKey {
	if !l.isExhausted && l.isStart {
		return &mKey{ValueOf(l.iter.Key())}
	}
	return nil
}
func (l *MapIter) Value() MValue {
	if !l.isExhausted && l.isStart {
		return &mValue{ValueOf(l.iter.Value())}
	}
	return nil
}

// input nil == r.*MapIter.Reset(r.Value{})
func (l *MapIter) Reset(v Map) {
	if v == nil {
		l.iter.Reset(r.Value{})
	} else {
		l.iter.Reset(v.Value())
		l.isStart = false
		l.isExhausted = false
	}
}

// ! >>>>>>>>>>>>>> MapSetter  <<<<<<<<<<<<<<<
type MapSetter = *vMapSetter
type vMapSetter struct {
	*vMap
	*vSetter
}

func (v MapSetter) init(p Pointer) ValueSetter {
	vMap := p.Elem().(Map)
	v = &vMapSetter{vMap, &vSetter{&vMap.v}}
	return v
}
func (v MapSetter) kind() r.Kind { return r.Map }

// func (v MapS)
