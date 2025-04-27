package gostd

import (
	"crypto/hmac"
	"crypto/md5"
	"io"
	"io/fs"
	"os"
	"testing"
)

/*
! New 返回的 MD5 Hash 新实例。它实现了 encoding.BinaryMarshaler 和 encoding.BinaryUnmarshaler 来封送和解包哈希的内部状态
! Sum 返回 data 的 MD5 校验和
*/

func TestMD5(t *testing.T) {
	h := md5.New()
	io.WriteString(h, "The fog is getting thicker!")
	io.WriteString(h, "And Leon's getting laaarger!")
	logfln("%x", h.Sum(nil))
	// e2c569be17396eca2a2e3c11578123ed
}

func TestMD5Files(t *testing.T) {
	// fs.
	fsys := os.DirFS(".")
	var fn func(t *testing.T, f io.Reader) (sum []byte)
	walkDir := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if err != io.EOF {
				return err
			}
		}
		if !d.IsDir() {
			f, err := fsys.Open(path)
			if err != nil {
				return err
			}
			logfln("%16s  >>> %x", d.Name(), fn(t, f))
			f.Close()
		}
		return nil
	}
	log(">>>>> md5 <<<<<")
	fn = md5File
	fs.WalkDir(fsys, ".", walkDir)

	log(">>>>> hmach md5 <<<<<")
	fn = hmach_md5File
	fs.WalkDir(fsys, ".", walkDir)

	/*
		>>>>> md5 <<<<<
		     aes_test.go  >>> 0f5ef0591edc87a6974bc7f1d770dfb5
		  cipher_test.go  >>> f104d1e6882ae7aed6a9dc1163f1774f
		         code.go  >>> 48c35984be9e7f94c0db7597ab1e3913
		  crypto_test.go  >>> 197d76d665c4db8a6fcb9ea795290741
		     des_test.go  >>> 373db93d38e699c9075834c944f53675
		    ecdh_test.go  >>> 7ac35763e70130656cacef3418a1793b
		   ecdsa_test.go  >>> a4845eca66d457137c0dac26386c93e8
		 ed25519_test.go  >>> c1bc1c926ba34377ad20629ac79eff7b
		elliptic_test.go  >>> 48ee6c629f7ad9b00fc413507c782ff0
		    hmac_test.go  >>> c778747617a3d4e557276c5795eed7a2
		     md5_test.go  >>> 6b7d05f6f14840f4cb949a5ef5205d9d
		>>>>> hmach md5 <<<<<
		     aes_test.go  >>> f13f9f656a47ab470afb7e4e5376b48e
		  cipher_test.go  >>> b8a42a001b50ca173d3de880846ca6c5
		         code.go  >>> 1ebd3efdf4f233ad3437b0ed5987ad8b
		  crypto_test.go  >>> d8d907a22d82e9634cffe46293fc3f81
		     des_test.go  >>> af94ba21e6a0196b491ffe51d4416e65
		    ecdh_test.go  >>> 4ef028a71581d0520b6d85da352764a7
		   ecdsa_test.go  >>> 3df01dce631ff53e3783334f4b8f6515
		 ed25519_test.go  >>> 814851cac63b91444d4fdfebc7738a12
		elliptic_test.go  >>> 88f623c35335eb45d619a8245ea4010e
		    hmac_test.go  >>> a4779651b2eb96272015d24fdd90e413
		     md5_test.go  >>> cc98b2494b771e557abad7ffb57c439e
	*/
}

func hmach_md5File(t *testing.T, f io.Reader) (sum []byte) {
	hmach := hmac.New(md5.New, []byte("Hello World"))
	if _, err := io.Copy(hmach, f); err != nil {
		t.Fatal(err)
	}
	return hmach.Sum(nil)
}

func md5File(t *testing.T, f io.Reader) []byte {
	d, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	sum := md5.Sum(d)
	return sum[:]
}
