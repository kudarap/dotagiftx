package hash

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"
)

// Sha1 returns a sha1 hashed string.
func Sha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s + Salt))

	return hex.EncodeToString(h.Sum(nil))
}

// GenerateSha1 returns a generated sha1 hash.
func GenerateSha1() string {
	t := time.Now().UnixNano()
	return Sha1(fmt.Sprintf("%d %s", t, Salt))
}
