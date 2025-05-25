package dgx

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type ChatHash string

type Chat struct {
	ID        string
	Users     []User
	Item      Item
	MarketID  string // reference from point of transaction
	Market    Market
	Messages  []Message
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Signature generates full conversation hash for anti-tampering protection.
func (c Chat) Signature(secret []byte) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write(b)
	h.Write(secret)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (c Chat) verify(secret []byte) error {
	return nil
}

func ChatVerify(secret, file []byte) bool {
	// strip first line to extract signature

	// parse file contents to chat struct

	return false
}

type Message struct {
	UserID    string
	SteamID   string
	Content   string
	Timestamp time.Time
}
