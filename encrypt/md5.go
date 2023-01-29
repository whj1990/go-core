package encrypt

import (
	"crypto/md5"
	"fmt"
	"io"
)

func Md5Encrypt(str string) string {
	w := md5.New()
	_, _ = io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}
