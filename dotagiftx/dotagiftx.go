package dotagiftx

import (
	"fmt"
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
	Active     string `json:"active"`

	started    time.Time
	timeFormat string
}

// NewVersion returns a formatted version details.
func NewVersion(prod bool, tag, commit, built string) *Version {
	v := &Version{
		Production: prod,
		Tag:        tag,
		Commit:     commit,
		Built:      built,

		started:    time.Now(),
		timeFormat: "Mon Jan 2 15:04:05 -0700 MST 2006",
	}
	v.formatBuiltDate()
	v.formatTag()
	return v
}

// formatBuiltDate format timestamp to human friendly dates.
func (v *Version) formatBuiltDate() {
	if strings.TrimSpace(v.Built) == "" {
		return
	}

	i, _ := strconv.ParseInt(v.Built, 10, 64)
	v.Built = time.Unix(i, 0).Format(v.timeFormat)
}

func (v *Version) SetUptime() *Version {
	v.Active = fmt.Sprintf(
		"Since %s; %s",
		v.started.Format(v.timeFormat),
		time.Since(v.started).Truncate(time.Second).String(),
	)
	return v
}

func (v *Version) formatTag() {
	if v.Production {
		return
	}

	v.Tag += "-dev"
}
