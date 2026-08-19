package dotagiftx

import "errors"

type (
	// FindOpts represents find options.
	FindOpts struct {
		Keyword       string
		KeywordFields []string
		Filter        any
		UserID        string
		Sort          string
		Desc          bool
		Page          int
		Limit         int
		Fields        []string
		WithMeta      bool
		// Advance options
		IndexSorting bool // Use for sorting indexed field.
		IndexKey     string
	}

	// FindMetadata represents find metadata.
	FindMetadata struct {
		ResultCount int
		TotalCount  int
	}
)

func (o FindOpts) validate() error {
	if o.Page < 0 || o.Limit < 0 {
		return errors.New("invalid page or limit")
	}
	return nil
}
