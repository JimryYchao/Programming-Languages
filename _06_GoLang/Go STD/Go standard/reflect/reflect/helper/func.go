package helper

import (
	r "reflect"
)

//! >>>>>>>>>>>>>> Func <<<<<<<<<<<<<<

type Func = *vFunc

type vFunc struct {
	*vCommon
}

func (v Func) valueof(rv r.Value) Value {
	v = &vFunc{newValue(rv)}
	return v
}
func (v Func) Kind() r.Kind { return r.Chan }
func (v Func) FuncType() FuncType {
	t, _ := TypeTo(v.Type()).FuncType()
	return t
}
