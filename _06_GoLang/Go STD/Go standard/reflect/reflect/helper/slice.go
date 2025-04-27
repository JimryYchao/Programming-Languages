package helper

import (
	r "reflect"
)

//! >>>>>>>>>>>>>> Slice <<<<<<<<<<<<<<

type Slice = *vSlice

type vSlice struct {
	*vCommon
}

func (v Slice) valueof(rv r.Value) Value {
	v = &vSlice{newValue(rv)}
	return v
}
func (v Slice) Kind() r.Kind { return r.Slice }
func (v Slice) SliceType() SliceType {
	t, _ := TypeTo(v.Type()).SliceType()
	return t
}

func (v Slice) Cap() int { return v.v.Cap() }
func (v Slice) Len() int { return v.v.Len() }
func (v Slice) Clear()   { v.v.Clear() }
func (v Slice) Index(i int) (Value, error) {
	if i < 0 || i >= v.v.Len() {
		return nil, newErr("Slice.Index", ErrOutOfRange)
	}
	return valueFrom(v.v.Index(i)), nil
}
func (s Slice) Slice(i, j int) (Slice, error) {
	if i < 0 || j < i || j > s.Len() {
		return nil, newErrorf("slice index out of bounds")
	}
	r := s.v.Slice(i, j)
	return &vSlice{newValue(r)}, nil
}

// ! >>>>>>>>>>>>>> SliceSetter <<<<<<<<<<<<<
type SliceSetter = *vSliceSetter

type vSliceSetter struct {
	*vSlice
	*vSetter
}

func (v SliceSetter) init(p Pointer) ValueSetter {
	vp := p.Elem().(Slice)
	v = &vSliceSetter{vp, &vSetter{&vp.v}}
	return v
}
func (v SliceSetter) kind() r.Kind { return r.Slice }

func (v SliceSetter) SetLen(n int) error {
	if uint(n) > uint(v.Cap()) {
		return newErrorf("slice length(%d) out of range in SetLen(%d)", n, v.Cap())
	}
	v.v.SetLen(n)
	return nil
}

// func (s Slice) String() string {
// 	sb := strings.Builder{}
// 	vs := make([]string, s.Len())
// 	for i := range len(vs) {
// 		si := s.Index(i)
// 		vs[i] = fmt.Sprint(si)
// 	}
// 	sb.WriteString("[")

// 	for i := range len(vs) - 1 {
// 		sb.WriteString(vs[i] + ", ")
// 	}
// 	sb.WriteString(vs[len(vs)-1] + "]")
// 	return sb.String()
// }

// func (s Slice) Interface() interface{} {
// 	return s.v.Interface()
// }

// func (s Slice) Slice3(i, j, k int) (Slice, error) {
// 	if i < 0 || j < i || k < j || k > s.Cap() {
// 		return nil, newErr("slice index out of bounds")
// 	}
// 	r := s.v.Slice3(i, j, k)
// 	return &vSlice{newValue(r)}, nil
// }

// 仅在 SliceFrom(&slice) 有效
// func (s Slice) Grow(n int) error {
// 	if s.p != nil {
// 		s.v.Grow(n)
// 		return nil
// 	} else {
// 		return newErr("s is not settable")
// 	}
// }

// func (s Slice) SetCap(n int) error {
// 	if n < s.Len() || n > s.Cap() {
// 		return newErr("slice capacity out of range in SetCap")
// 	}
// 	s.v.SetCap(n)
// 	return nil
// }
// func (s Slice) IndexValue(i int) *r.Value {
// 	if uint(i) >= uint(s.Len()) {
// 		return nil
// 	}
// 	r := s.v.Index(i)
// 	return &r
// }

// func (s Slice) Index(i int) any {
// 	return s.IndexValue(i).Interface()
// }

// func (s Slice) SetIndexValue(i int, v r.Value) bool {
// 	if s.v.Type().Elem().ConvertibleTo(v.Type()) {
// 		s.v.Index(i).Set(v)
// 		return true
// 	}
// 	return false
// }
// func (s Slice) SetIndex(i int, v any) bool {
// 	return s.SetIndexValue(i, r.ValueOf(v))
// }

// func (s Slice) Value() r.Value {
// 	return *s.v
// }

// func (s Slice) Append(x ...any) (Slice, error) {
// 	if len(x) == 0 {
// 		return s, nil
// 	}

// 	if !r.TypeOf(x[0]).AssignableTo(s.Type().Elem()) {
// 		return nil, newErr(fmt.Sprintf("the type of x.Elem(%s) is not s.Elem: %s", x, s.ElemString()))
// 	}
// 	xs := make([]r.Value, len(x))
// 	for i, v := range x {
// 		xs[i] = r.ValueOf(v)
// 	}
// 	return s.AppendValue(xs...)
// }
// func (s Slice) AppendSlice(x r.Value) (Slice, error) {
// 	if x.Kind() != r.Slice {
// 		return nil, newErr("x is not a slice")
// 	}
// 	if !x.Type().Elem().AssignableTo(s.Type().Elem()) {
// 		return nil, newErr(fmt.Sprintf("the type of x.Elem(%s) is not s.Elem: %s", x.Type().Elem(), s.ElemString()))
// 	}
// 	nv := r.AppendSlice(s.Value(), x)
// 	return &vSlice{&nv, nil}, nil
// }

// func (s Slice) AppendValue(x ...r.Value) (Slice, error) {
// 	et := s.Type().Elem()
// 	var rt r.Value = *s.v
// 	for _, v := range x {
// 		if v.Type() != et {
// 			return nil, newErr(fmt.Sprintf("the type of %v is not %s", v, et))
// 		}
// 		rt = r.Append(rt, v)
// 	}
// 	return &vSlice{&rt, nil}, nil
// }

// func SliceFrom(s interface{}) (Slice, error) {
// 	v := r.ValueOf(s)
// 	k := v.Kind()
// 	if v.Kind() == r.Pointer {
// 		if v.Elem().Kind() == r.Slice {
// 			sv := v.Elem()
// 			return &vSlice{&sv, &v}, nil
// 		} else {
// 			return nil, fmt.Errorf("s is not a pointer to slice")
// 		}
// 	}
// 	if k != r.Slice {
// 		return nil, fmt.Errorf("s is not a slice: kind=%s", v.Kind())
// 	}
// 	return &vSlice{&v, nil}, nil
// }

// func SliceFromValue(v r.Value) (Slice, error) {
// 	if v.Kind() != r.Slice {
// 		return nil, fmt.Errorf("s is not a slice: kind=%s", v.Kind())
// 	}
// 	return &vSlice{&v, nil}, nil
// }
