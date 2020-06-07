package hash

import "github.com/speps/go-hashids"

// EncodeID encodes integer inputs to hash id string.
func EncodeID(len int, input ...int) string {
	h := hashid(len)
	e, _ := h.Encode(input)
	return e
}

// DecodeID decodes hash id into integers.
func DecodeID(len int, hash string) []int {
	h := hashid(len)
	d, _ := h.DecodeWithError(hash)
	return d
}

func hashid(len int) *hashids.HashID {
	if len == 0 {
		len = 10
	}

	hd := hashids.NewData()
	hd.Salt = Salt
	hd.MinLength = len
	h, _ := hashids.NewWithData(hd)
	return h
}
