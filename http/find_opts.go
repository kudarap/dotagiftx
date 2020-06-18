package http

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/schema"
	"github.com/kudarap/dota2giftables/core"
)

const defaultPageLimit = 10

func findOptsFromURL(u *url.URL, filter interface{}) (core.FindOpts, error) {
	opts := core.FindOpts{}
	get := u.Query().Get

	// Set pagination.
	opts.Page, _ = strconv.Atoi(get("page"))
	opts.Limit, _ = strconv.Atoi(get("limit"))
	if opts.Limit == 0 {
		opts.Limit = defaultPageLimit
	}

	// Sets sort.
	opts.Sort, opts.Desc = parseSort(get("sort"))

	// Set filter.
	if err := findOptsFilter(u, filter); err != nil {
		return core.FindOpts{}, err
	}
	opts.Filter = filter
	opts.WithMeta = true

	return opts, nil
}

const sortDescSuffix = ":desc"

func parseSort(sortStr string) (field string, isDesc bool) {
	// Get sort field.
	s := strings.Split(sortStr, ":")
	if len(s) == 0 {
		return
	}
	field = s[0]

	// Detect sorting order.
	if strings.HasSuffix(sortStr, sortDescSuffix) {
		isDesc = true
		return
	}

	return
}

const defaultFilterTag = "json"

func findOptsFilter(u *url.URL, filter interface{}) error {
	// Sets search filters.
	d := schema.NewDecoder()
	d.SetAliasTag(defaultFilterTag)
	d.IgnoreUnknownKeys(true)
	return d.Decode(filter, u.Query())
}
