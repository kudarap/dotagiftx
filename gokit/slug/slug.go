package slug

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile("[^a-z0-9]+")

// Make creates a URL friendly string base on input.
func Make(s string) string {
	return strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
}
