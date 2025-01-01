package helper

import r "reflect"

type String = *vString
type vString struct {
	*vCommon
}

func (s String) valueof(rv r.Value) Value {
	s = &vString{newValue(rv)}
	return s
}
func (s String) Kind() r.Kind { return r.String }
func (s String) StringType() TypeCommon {
	return newTCommon(s.v.Type())
}

func (s String) Len() int { return s.v.Len() }
func (s String) Index(i int) (byte, error) {
	if i < 0 || i >= s.v.Len() {
		return 0, newErr("String.Index", ErrOutOfRange)
	}
	return s.v.Index(i).Interface().(byte), nil
}

func (s String) SetString(x string) bool {
	if s.CanSet() {
		s.v.SetString(x)
		return true
	}
	return false
}

func (s String) String() string { return s.v.String() }
