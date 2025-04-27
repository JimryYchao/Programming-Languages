package gostd

import (
	"crypto/sha512"
	"testing"
)

/*
! New, New384, New512_224, New512_256
! Sum512, Sum384, Sum512_224, Sum512_256
*/

func TestSha512(t *testing.T) {
	data := []byte("His money is twice tainted: 'taint yours and 'taint mine.")
	logfln("sha512 : %x", sha512.Sum512(data))
	logfln("sha384 : %x", sha512.Sum384(data))
	logfln("sha512_224 : %x", sha512.Sum512_224(data))
	logfln("sha512_256 : %x", sha512.Sum512_256(data))
}
