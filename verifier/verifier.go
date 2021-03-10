package verifier

type VerifyStatus uint

const (
	VerifyStatusPrivate = iota
	VerifyStatusItem
	VerifyStatusSeller
)

type Asset struct {
}
