package helper_test

import (
	"gostd/reflect/helper"
	"testing"
)

func TestStringValue(t *testing.T) {
	s, ok := helper.ValueSpecifyOf[helper.String]("Hello")
	if ok {
		b, _ := s.Index(1)
		logf("%c", b)
	}

	st := s.StringType()

	log(st.Kind())

	i, ok := helper.ValueSpecifyOf[helper.Int](12)
	if t, ok := i.Type().(helper.TypeCommon); ok {
		log(t.Kind())
		log(helper.Zero(t))
		log(helper.Zero(helper.TypeFor[helper.Value]()))
	}
}
