package version

import (
	"strconv"
	"strings"
	"time"
)

// Version represents application version.
type Version struct {
	Production bool   `json:"production"`
	Tag        string `json:"version"`
	Commit     string `json:"hash"`
	Built      string `json:"built"`
}

// New returns a formatted version details.
func New(prod bool, tag, commit, built string) *Version {
	v := &Version{
		prod,
		tag,
		commit,
		built,
	}
	v.formatBuiltDate()
	v.formatTag()
	return v
}

// formatBuiltDate formats timestamp to human friendly dates.
func (v *Version) formatBuiltDate() {
	if strings.TrimSpace(v.Built) == "" {
		return
	}

	i, _ := strconv.ParseInt(v.Built, 10, 64)
	v.Built = time.Unix(i, 0).Format("Mon Jan 2 15:04:05 -0700 MST 2006")
}

func (v *Version) formatTag() {
	if v.Production {
		return
	}

	v.Tag += "-dev"
}
