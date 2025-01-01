package helper

import (
	r "reflect"
)

//! >>>>>>>>>>>>>> Chan <<<<<<<<<<<<<<

type Chan = *vChan
type ChanB = *vChanB

// type ChanR = *vChanR
// type ChanS = *vChanS

type vChan struct {
	*vCommon
	dir r.ChanDir
}

func (v Chan) valueof(rv r.Value) Value {
	v = &vChan{newValue(rv), rv.Type().ChanDir()}
	return v
}
func (v Chan) Kind() r.Kind { return r.Chan }
func (v Chan) ChanType() ChanType {
	t, _ := TypeTo(v.Type()).ChanType()
	return t
}

func (v Chan) Cap() int { return v.v.Cap() }
func (v Chan) Len() int { return v.v.Len() }

func (v Chan) IsNil() bool { return v.v.IsNil() }

// 补充
func (v Chan) ChanDir() r.ChanDir { return v.dir }

// both chan
func (v Chan) ToChanB() (ChanB, bool) {
	if v.dir == r.BothDir {
		return &vChanB{v}, true
	}
	return nil, false
}

type vChanB struct {
	ChanValue Chan
}
