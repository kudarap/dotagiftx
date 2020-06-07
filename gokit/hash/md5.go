package hash

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

// MD5 returns a md5 hashed string.
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s + Salt))

	return hex.EncodeToString(h.Sum(nil))
}

// GenerateMD5 returns a generated md5 hash.
func GenerateMD5() string {
	t := time.Now().UnixNano()
	return MD5(fmt.Sprintf("%d %s", t, Salt))
}
