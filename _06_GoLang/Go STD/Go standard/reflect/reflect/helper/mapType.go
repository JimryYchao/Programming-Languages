package helper

import (
	"fmt"
	r "reflect"
)

//! >>>>>>>>>>>> MapType <<<<<<<<<<<<

type MapType = *mapType

type mapType struct {
	*tCommon
}

func (t MapType) typeof(tp r.Type) Type {
	t = &mapType{newTCommon(tp)}
	return t
}

func (MapType) Kind() r.Kind { return r.Map }
func (t MapType) Elem() Type { return typeFrom(t.t.Elem()) }
func (t MapType) Key() Type  { return typeFrom(t.t.Key()) }

// MapOf
func MapOf(key r.Type, elem r.Type) (MapType, error) {
	if key == nil || elem == nil {
		return nil, newErr("MapOf", ErrArgNil)
	}
	if !key.Comparable() {
		return nil, newErr("MapOf", fmt.Errorf("invalid key type(%s)", key))
	}

	mtp := r.MapOf(key, elem)
	return &mapType{newTCommon(mtp)}, nil
}

func MapFor[K comparable, V any]() MapType {
	m, _ := MapOf(r.TypeFor[K](), r.TypeFor[V]())
	return m
}

// new
