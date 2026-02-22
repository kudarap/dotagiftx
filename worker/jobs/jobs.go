package jobs

type cacheRemover interface {
	BulkDel(keyPrefix string) error
}
