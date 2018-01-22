package modle

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
)

func Md5(str string) string {
	m := md5.New()
	io.WriteString(m, str)
	return strings.ToUpper(fmt.Sprintf("%x", m.Sum(nil)))
}
