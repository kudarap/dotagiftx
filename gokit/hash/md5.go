package hash

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 returns a md5 hashed string.
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s + Salt))

	return hex.EncodeToString(h.Sum(nil))
}
