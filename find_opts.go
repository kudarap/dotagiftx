package dgx

type (
	// FindOpts represents find options.
	FindOpts struct {
		Keyword       string
		KeywordFields []string
		Filter        interface{}
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
