package rethink

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/kudarap/dotagiftx/core"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type findOpts core.FindOpts

func newFindOptsQuery(q r.Term, o core.FindOpts) r.Term {
	//return findOpts(o).parseOpts(q, nil)
	return baseFindOptsQuery(q, o, nil)
}

func baseFindOptsQuery(q r.Term, o core.FindOpts, hookFn func(r.Term) r.Term) r.Term {
	return findOpts(o).parseOpts(q, hookFn)
}

func (o findOpts) parseOpts(q r.Term, hookFn func(r.Term) r.Term) r.Term {
	if o.IndexSorting && o.Sort != "" {
		q = q.OrderBy(r.OrderByOpts{Index: o.parseOrder()})
	}

	if hookFn != nil {
		q = hookFn(q)
	}

	if strings.TrimSpace(o.Keyword) != "" {
		q = q.Filter(o.parseKeyword())
	}

	if o.Filter != nil {
		q = q.Filter(o.parseFilter())
	}

	if o.UserID != "" {
		q = q.Filter(o.setUserScope())
	}

	if !o.IndexSorting && o.Sort != "" {
		q = q.OrderBy(o.parseOrder())
	}

	if o.Limit != 0 {
		q = q.Slice(o.parseSlice())
	}

	if o.Fields != nil {
		q = q.Pluck(o.Fields)
	}

	return q
}

func (o findOpts) parseKeyword() interface{} {
	if len(o.KeywordFields) == 0 {
		return nil
	}

	return func(t r.Term) r.Term {
		// Concatenate values of search fields to create a fake index.
		searchText := t.Field(o.KeywordFields[0])
		for _, ff := range o.KeywordFields[1:] {
			searchText = searchText.Add(" ", t.Field(ff))
		}

		// Matches that contains the keywords non case sensitive.
		q := searchText
		for _, ww := range strings.Split(normalizeKeyword(o.Keyword), " ") {
			q = q.And(searchText.Match(fmt.Sprintf("(?i)%s", ww)))
		}

		return q
	}
}

// normalizeKeyword handles special case for the word "Collector's" with apostrophe.
func normalizeKeyword(keyword string) string {
	s := strings.ToLower(keyword)

	// Special case for the word "Collector's" with apostrophe.
	if strings.Contains(s, "collectors") {
		s = strings.ReplaceAll(s, "collectors", "collector's")
	}

	return s
}

func (o findOpts) parseFilter() map[string]interface{} {
	if o.Filter == nil {
		return map[string]interface{}{}
	}

	structs.DefaultTagName = tagName
	return structs.New(o.Filter).Map()
}

func (o findOpts) parseOrder() interface{} {
	if o.Desc {
		return r.Desc(o.Sort)
	}

	return o.Sort
}

func (o findOpts) parseSlice() (start int, end int) {
	if o.Page < 1 {
		o.Page = 1
	}
	o.Page--

	start = o.Page * o.Limit
	end = start + o.Limit
	return
}

func (o findOpts) setUserScope() map[string]interface{} {
	return map[string]interface{}{
		"user_id": o.UserID,
	}
}
