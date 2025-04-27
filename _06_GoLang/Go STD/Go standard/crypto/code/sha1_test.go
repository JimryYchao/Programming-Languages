package gostd

import (
	"crypto/sha1"
	"io"
	"os"
	"testing"
)

// ! New & Sum
func TestSha1Sum(t *testing.T) {
	h := sha1.New()
	io.WriteString(h, "His money is twice tainted:")
	io.WriteString(h, " 'taint yours and 'taint mine.")
	logfln("%x", h.Sum(nil))

	logfln("%x", sha1.Sum([]byte("His money is twice tainted: 'taint yours and 'taint mine.")))
}

func TestSha1File(t *testing.T) {
	f, err := os.Open("sha1_test.go")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		t.Fatal(err)
	}
	logfln("%x", h.Sum(nil))
}
