package rethink

import (
	"fmt"
	"log"

	"github.com/kudarap/dotagiftx/core"
)

const (
	tableQueue = "queue"
)

type queueStorage struct {
	db *Client
}

func NewQueue(c *Client) *queueStorage {
	if err := c.autoMigrate(tableQueue); err != nil {
		log.Fatalf("could not create %s table: %s", tableTrack, err)
	}

	if err := c.autoIndex(tableQueue, core.Track{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableMarket, err)
	}

	return &queueStorage{c}
}

func (q *queueStorage) VerifyDelivery(marketID string) {
	fmt.Println("VerifyDelivery", marketID)
}

func (q *queueStorage) VerifyInventory(userID string) {
	fmt.Println("VerifyInventory", userID)
}
