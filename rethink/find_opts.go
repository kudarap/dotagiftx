package rethink

import (
	"github.com/fatih/structs"
	"github.com/kudarap/dota2giftables/core"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type findOpts core.FindOpts

func newFindOptsQuery(q r.Term, o core.FindOpts) r.Term {
	return findOpts(o).parseOpts(q)
}

func (o findOpts) parseOpts(q r.Term) r.Term {
	if o.Filter != nil {
		q = q.Filter(o.parseFilter())
	}

	if o.UserID != "" {
		q = q.Filter(o.setUserScope())
	}

	if o.Sort != "" {
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
