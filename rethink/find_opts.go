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
	return findOpts(o).parseOpts(q, nil)
}

func newCatalogFindOptsQuery(q r.Term, o core.FindOpts, hookFn func(r.Term) r.Term) r.Term {
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
		tags := t.Field(o.KeywordFields[0])
		for _, kf := range o.KeywordFields[1:] {
			tags = tags.Add(" ", t.Field(kf))
		}

		kws := strings.Split(o.Keyword, " ")
		q := tags.Match(fmt.Sprintf("(?i)%s", kws[0]))
		for _, kw := range kws[1:] {
			q = q.And(tags.Match(fmt.Sprintf("(?i)%s", kw)))
		}

		// Matches that contains the keyword non case sensitive.
		return q
	}
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
