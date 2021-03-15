package jobs

import (
	"fmt"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/worker"
)

// Dispatcher represents the handling of custom jobs.
type Dispatcher struct {
	worker *worker.Worker
	// Jobs service dependencies
	deliverySvc  core.DeliveryService
	inventorySvc core.InventoryService
	marketSvc    core.MarketService
	logger       log.Logger
}

// NewDispatcher returns an instance dispatcher.
func NewDispatcher(worker *worker.Worker,
	deliverySvc core.DeliveryService,
	inventorySvc core.InventoryService,
	marketSvc core.MarketService,
	logger log.Logger,
) *Dispatcher {
	if worker == nil {
		panic("worker is nul")
	}

	return &Dispatcher{
		worker,
		deliverySvc,
		inventorySvc,
		marketSvc,
		logger,
	}
}

// VerifyDelivery creates a job to verify a delivery
// and queue them to worker.
//
// Customized existing VerifyDelivery job to make it run once job.
func (d *Dispatcher) VerifyDelivery(marketID string) {
	job := NewVerifyDelivery(d.deliverySvc, d.marketSvc, d.logger)
	job.name = fmt.Sprintf("%s_%s", job.name, marketID)
	job.interval = 0 // makes the job run-once.
	job.filter = core.Market{ID: marketID}
	d.worker.AddJob(job)
}

// VerifyInventory creates a job to verify a inventory
// and queue them to worker.
//
// Customized existing VerifyInventory job to make it run once job.
func (d *Dispatcher) VerifyInventory(userID string) {
	job := NewVerifyInventory(d.inventorySvc, d.marketSvc, d.logger)
	job.name = fmt.Sprintf("%s_%s", job.name, userID)
	job.interval = 0 // makes the job run-once.
	job.filter.UserID = userID
	d.worker.AddJob(job)
}
