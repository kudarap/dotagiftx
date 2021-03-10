package verifier

type VerifyStatus string

const (
	VerifyStatusError   VerifyStatus = "error"
	VerifyStatusPrivate VerifyStatus = "private"
	VerifyStatusNoHit   VerifyStatus = "no-hit"
	VerifyStatusItem    VerifyStatus = "item"
	VerifyStatusSeller  VerifyStatus = "seller"
)

type Asset struct {
}
