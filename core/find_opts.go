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

	// Metadata represents find metadata.
	Metadata struct {
		ResultCount int
		TotalCount  int
	}
)
