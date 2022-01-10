// +build amd64

package image

import (
	"github.com/libgoost/encoding-base64"
)

func EncodeToString(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}
