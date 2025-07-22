package dotagiftx

const (
	// StorageUncaughtErr storage error type for un-handled errors.
	StorageUncaughtErr Errors = iota + storageErrorIndex
	// StorageMergeErr storage object merge error.
	StorageMergeErr
)

// init sets error text definition.
func init() {
	appErrorText[StorageUncaughtErr] = "un-handled storage error"
	appErrorText[StorageMergeErr] = "object merge error"
}
