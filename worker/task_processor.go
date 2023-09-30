package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/steaminv"
	"github.com/kudarap/dotagiftx/verified"
)

type TaskProcessor struct {
	queue taskQueue
	rate  time.Duration

	inventorySvc core.InventoryService
	deliverySvc  core.DeliveryService
}

func NewTaskProcessor(
	rate time.Duration,
	queue taskQueue,
	inventorySvc core.InventoryService,
	deliverySvc core.DeliveryService,
) *TaskProcessor {
	return &TaskProcessor{
		queue:        queue,
		rate:         rate,
		inventorySvc: inventorySvc,
		deliverySvc:  deliverySvc,
	}
}

func (p *TaskProcessor) Run(wg *sync.WaitGroup) {
	ctx := context.Background()
	for {
		time.Sleep(p.rate)

		start := time.Now()

		t, err := p.queue.Get(ctx)
		if err != nil {
			log.Printf("ERR! could not get task from queue: %s", err)
			continue
		}
		if t == nil {
			continue
		}

		task := *t
		wg.Add(1)

		task.Status = core.TaskStatusProcessing
		if err = p.queue.Update(ctx, task); err != nil {
			log.Printf("ERR! could not process task: %s", err)
			wg.Done()
			continue
		}
		log.Println("task get", task.ID, task.Type, task.Priority)

		var run func(context.Context, interface{}) error
		switch task.Type {
		case core.TaskTypeVerifyInventory:
			run = p.taskVerifyInventory
		case core.TaskTypeVerifyDelivery:
			run = p.taskVerifyDelivery
		}

		log.Println("task processing...", task.ID, task.Type)
		err = run(ctx, task.Payload)
		task.ElapsedMs = time.Since(start).Milliseconds()
		if err != nil {
			log.Printf("ERR! running tasks: %s %s", task.Type, err)
			task.Status = core.TaskStatusError
			task.Note = fmt.Sprintf("err: %s", err)
			if err = p.queue.Update(ctx, task); err != nil {
				log.Printf("ERR! could not run task: %s", err)
			}
			wg.Done()
			continue
		}

		task.Status = core.TaskStatusDone
		log.Println("task done!", task.ID, time.Duration(task.ElapsedMs)*time.Millisecond)
		if err = p.queue.Update(ctx, task); err != nil {
			log.Printf("ERR! could not update task: %s", err)
		}
		wg.Done()
	}
}

func (p *TaskProcessor) taskVerifyInventory(ctx context.Context, data interface{}) error {
	var market core.Market
	if err := marshallTaskPayload(data, &market); err != nil {
		return err
	}

	if market.User == nil || market.Item == nil {
		return fmt.Errorf("skipped process! missing data user:%#v item:%#v", market.User, market.Item)
	}

	source := steaminv.InventoryAssetWithCache
	status, assets, err := verified.Inventory(source, market.User.SteamID, market.Item.Name)
	if err != nil {
		return err
	}

	err = p.inventorySvc.Set(ctx, &core.Inventory{
		MarketID: market.ID,
		Status:   status,
		Assets:   assets,
	})
	return nil
}

func (p *TaskProcessor) taskVerifyDelivery(ctx context.Context, data interface{}) error {
	var mkt core.Market
	if err := marshallTaskPayload(data, &mkt); err != nil {
		return err
	}

	if mkt.User == nil || mkt.Item == nil {
		return fmt.Errorf("skipped process! missing data user:%#v item:%#v", mkt.User, mkt.Item)
	}

	src := steaminv.InventoryAsset
	status, assets, err := verified.Delivery(src, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
	if err != nil {
		return err
	}

	err = p.deliverySvc.Set(ctx, &core.Delivery{
		MarketID: mkt.ID,
		Status:   status,
		Assets:   assets,
	})
	return err
}

type taskQueue interface {
	Get(ctx context.Context) (*core.Task, error)
	Update(ctx context.Context, t core.Task) error
}

func marshallTaskPayload(in, out interface{}) error {
	raw, ok := in.(map[string]interface{})
	if !ok {
		return fmt.Errorf("un-supported payload")
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, out); err != nil {
		return err
	}
	return nil
}
