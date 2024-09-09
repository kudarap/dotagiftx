package rethink

import (
	"fmt"
	"log"
	"regexp"
	"time"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const tableSpan = "span"

type SpanStorage struct {
	db *Client
}

type span struct {
	Name      string    `db:"name,index"`
	ElapsedMs int64     `db:"elapsed_ms,index"`
	CreatedAt time.Time `db:"created_at,index"`
}

func NewSpan(c *Client) *SpanStorage {
	if err := c.autoMigrate(tableSpan); err != nil {
		log.Fatalf("could not create %s table: %s", tableSpan, err)
	}

	if err := c.autoIndex(tableSpan, span{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableSpan, err)
	}

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
	return spanUUID.ReplaceAllString(s, "<uuid>")
}
