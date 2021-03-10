package verifier

type VerifyStatus uint

const (
	VerifyStatusError = iota
	VerifyStatusPrivate
	VerifyStatusItem
	VerifyStatusSeller
)

type Asset struct {
}
