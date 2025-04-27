package helper

import (
	r "reflect"
)

//! >>>>>>>>>>>>>> Struct <<<<<<<<<<<<<<

type Struct = *vStruct

type vStruct struct {
	*vCommon
}

func (v Struct) valueof(rv r.Value) Value {
	v = &vStruct{newValue(rv)}
	return v
}
func (v Struct) Kind() r.Kind { return r.Struct }

func (v Struct) StructType() StructType {
	t, _ := TypeTo(v.Type()).StructType()
	return t
}
