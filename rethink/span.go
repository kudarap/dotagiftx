package rethink

import (
	"fmt"
	"regexp"
	"time"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type SpanStorage struct {
	db *Client
}

type span struct {
	Name      string    `db:"name"`
	ElapsedMs int64     `db:"elapsed_ms"`
	CreatedAt time.Time `db:"created_at"`
}

func NewSpan(c *Client) *SpanStorage {
	return &SpanStorage{c}
}

func (s *SpanStorage) Add(name string, elapsedMs int64, t time.Time) {
	name = spanCleanUUIDs(name)
	i := span{name, elapsedMs, t}
	_, err := s.db.insert(r.Table("span").Insert(i))
	if err != nil {
		fmt.Println("ERR SPAN", err)
	}
}

var spanUUID = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func spanCleanUUIDs(s string) string {
	return spanUUID.ReplaceAllString(s, "{UUID}")
}
