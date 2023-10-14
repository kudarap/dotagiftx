package rethink

import (
	"fmt"
	"time"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type spanStorage struct {
	db *Client
}

type span struct {
	Name      string    `db:"name"`
	ElapsedMs int64     `db:"elapsed_ms"`
	CreatedAt time.Time `db:"created_at"`
}

func NewSpan(c *Client) *spanStorage {
	return &spanStorage{c}
}

func (s *spanStorage) Add(name string, elapsedMs int64, t time.Time) {
	i := span{name, elapsedMs, t}
	_, err := s.db.insert(r.Table("span").Insert(i))
	if err != nil {
		fmt.Println("ERR SPAN", err)
	}
}
