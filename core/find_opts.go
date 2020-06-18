package core

type (
	// FindOpts represents find options.
	FindOpts struct {
		Filter   interface{}
		UserID   string
		Sort     string
		Desc     bool
		Page     int
		Limit    int
		Fields   []string
		WithMeta bool
	}

	// FindMetadata represents find metadata.
	FindMetadata struct {
		ResultCount int
		TotalCount  int
	}
)
