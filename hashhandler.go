package FileHasher

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func hashFile(f string) []byte {
	if b, err := os.Open(f); err != nil {
		panic(err)
	} else {
		hsh := sha256.New()
		if _, er := io.Copy(hsh, b); er != nil {
			panic(er)
		}
		h := []byte(fmt.Sprintf("%x", hsh.Sum(nil)))
		return h
	}
}
