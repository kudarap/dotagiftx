package jobs

import (
	"fmt"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/worker"
	"github.com/sirupsen/logrus"
)

// Dispatcher represents the handling of custom jobs.
type Dispatcher struct {
	worker *worker.Worker
	// Jobs service dependencies
	deliverySvc  core.DeliveryService
	inventorySvc core.InventoryService
	marketStg    core.MarketStorage
	logSvc       *logrus.Logger
}

// NewDispatcher returns an instance dispatcher.
func NewDispatcher(worker *worker.Worker,
	deliverySvc core.DeliveryService,
	inventorySvc core.InventoryService,
	marketStg core.MarketStorage,
	logSvc *logrus.Logger,
) *Dispatcher {
	return &Dispatcher{
		worker,
		deliverySvc,
		inventorySvc,
		marketStg,
		logSvc,
	}
}

// RegisterJobs add pre-defined jobs, mostly recurring one's.
func (d *Dispatcher) RegisterJobs() {
	d.worker.AddJob(NewRecheckInventory(
		d.inventorySvc,
		d.marketStg,
		log.WithPrefix(d.logSvc, "job_recheck_inventory"),
	))
	d.worker.AddJob(NewVerifyInventory(
		d.inventorySvc,
		d.marketStg,
		log.WithPrefix(d.logSvc, "job_verify_inventory"),
	))
	//d.worker.AddJob(NewVerifyDelivery(
	//	d.deliverySvc,
	//	d.marketStg,
	//	log.WithPrefix(d.logSvc, "job_verify_delivery"),
	//))
	d.worker.AddJob(NewGiftWrappedUpdate(
		d.deliverySvc,
		d.marketStg,
		log.WithPrefix(d.logSvc, "job_giftwrapped_update"),
	))
}

// VerifyDelivery creates a job to verify a delivery
// and queue them to worker.
//
// Customized existing VerifyDelivery job to make it run once job.
func (d *Dispatcher) VerifyDelivery(marketID string) {
	ctxLog := log.WithPrefix(d.logSvc, "dispatch_verify_inventory")
	job := NewVerifyDelivery(d.deliverySvc, d.marketStg, ctxLog)
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
	ctxLog := log.WithPrefix(d.logSvc, "dispatch_verify_inventory")
	job := NewVerifyInventory(d.inventorySvc, d.marketStg, ctxLog)
	job.name = fmt.Sprintf("%s_%s", job.name, userID)
	job.interval = 0 // makes the job run-once.
	job.filter.UserID = userID
	d.worker.AddJob(job)
}
