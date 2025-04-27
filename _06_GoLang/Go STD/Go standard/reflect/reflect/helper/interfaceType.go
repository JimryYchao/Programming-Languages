package helper

import (
	"errors"
	r "reflect"
)

//! >>>>>>>>>>>> InterfaceType <<<<<<<<<<<<

type InterfaceType = *vInterfaceType

type vInterfaceType struct {
	tCommonBase
}

func (t InterfaceType) typeof(tp r.Type) Type {
	it := vInterfaceType{}
	it.tCommonBase = newTCommon(tp)
	t = &it
	return t
}

func (InterfaceType) Kind() r.Kind { return r.Interface }

func InterfaceFor[T any]() (InterfaceType, error) {
	tp := r.TypeFor[T]()
	if tp.Kind() != r.Interface {
		return nil, newErr("InterfaceFor", errors.New("Type is not a interface"))
	}
	it := vInterfaceType{}
	it.tCommonBase = newTCommon(tp)
	return &it, nil
}
