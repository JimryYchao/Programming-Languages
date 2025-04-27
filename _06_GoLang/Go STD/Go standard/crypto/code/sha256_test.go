package gostd

import (
	"crypto/sha256"
	"io"
	"os"
	"testing"
)

/*
! New, Sum256
! New224, Sum224
*/

func TestSha256(t *testing.T) {
	h := sha256.New()
	io.WriteString(h, "His money is twice tainted:")
	io.WriteString(h, " 'taint yours and 'taint mine.")
	logfln("sha256 : %x", h.Sum(nil))

	logfln("sha224 : %x", sha256.Sum224([]byte("His money is twice tainted: 'taint yours and 'taint mine.")))
}

func TestSha256File(t *testing.T) {
	f, err := os.Open("sha256_test.go")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		t.Fatal(err)
	}
	logfln("%x", h.Sum(nil))
}
