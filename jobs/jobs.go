// This package should only contain a worker jobs.

package jobs

import (
	"fmt"
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/worker"
	"github.com/sirupsen/logrus"
)

// const defaultJobInterval = time.Hour * 24
const defaultJobInterval = time.Hour * 6

// Dispatcher represents the handling of custom jobs.
type Dispatcher struct {
	worker *worker.Worker
	// Jobs service dependencies
	deliverySvc  dgx.DeliveryService
	inventorySvc dgx.InventoryService
	deliveryStg  dgx.DeliveryStorage
	marketStg    dgx.MarketStorage
	catalogStg   dgx.CatalogStorage
	userStg      dgx.UserStorage
	cache        dgx.Cache
	logSvc       *logrus.Logger
}

// NewDispatcher returns an instance dispatcher.
func NewDispatcher(
	worker *worker.Worker,
	deliverySvc dgx.DeliveryService,
	inventorySvc dgx.InventoryService,
	deliveryStg dgx.DeliveryStorage,
	marketStg dgx.MarketStorage,
	catalogStg dgx.CatalogStorage,
	userStg dgx.UserStorage,
	cache dgx.Cache,
	logSvc *logrus.Logger,
) *Dispatcher {
	return &Dispatcher{
		worker,
		deliverySvc,
		inventorySvc,
		deliveryStg,
		marketStg,
		catalogStg,
		userStg,
		cache,
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
	d.worker.AddJob(NewVerifyDelivery(
		d.deliverySvc,
		d.marketStg,
		log.WithPrefix(d.logSvc, "job_verify_delivery"),
	))
	d.worker.AddJob(NewGiftWrappedUpdate(
		d.deliverySvc,
		d.deliveryStg,
		d.marketStg,
		log.WithPrefix(d.logSvc, "job_giftwrapped_update"),
	))
	d.worker.AddJob(NewRevalidateDelivery(
		d.deliverySvc,
		d.marketStg,
		log.WithPrefix(d.logSvc, "job_revalidate_delivery"),
	))
	d.worker.AddJob(NewExpiringSubscription(
		d.userStg,
		d.cache,
		log.WithPrefix(d.logSvc, "job_expiring_subscription"),
	))
	d.worker.AddJob(NewExpiringMarket(
		d.marketStg,
		d.catalogStg,
		d.cache,
		log.WithPrefix(d.logSvc, "job_expiring_market"),
	))
	d.worker.AddJob(NewSweepMarket(
		d.marketStg, log.WithPrefix(d.logSvc, "job_sweep_market"),
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
	job.filter = dgx.Market{ID: marketID}
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
