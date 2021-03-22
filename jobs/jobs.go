// This package should only contain a worker jobs.

package jobs

import (
	"math/rand"
	"time"
)

const defaultJobInterval = time.Hour * 24

func rest(maxSleep int) {
	rand.Seed(time.Now().Unix())

	t := time.Duration(rand.Intn(maxSleep-0) + 0)
	time.Sleep(time.Second * t)
}
