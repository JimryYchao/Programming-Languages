package helper

import r "reflect"

type Interface = *vInterface

type vInterface struct {
	*vCommon
}

func (v Interface) valueof(rv r.Value) Value {
	v = &vInterface{newValue(rv)}
	return v
}
func (v Interface) Kind() r.Kind { return r.Array }
func (v Interface) InterfaceType() InterfaceType {
	t, _ := TypeTo(v.Type()).InterfaceType()
	return t
}

func (v Interface) Elem() Value { return valueFrom(v.v.Elem()) }
