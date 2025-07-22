package jobs

// cache provides access to cache database.
type cache interface {
	BulkDel(keyPrefix string) error
}
