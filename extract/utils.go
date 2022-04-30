package extract

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
)

func GetMD5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func NullUnPadding(in []byte) []byte {
	return bytes.TrimRight(in, string([]byte{0}))
}
