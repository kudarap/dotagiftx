package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/verifying"
)

type TaskProcessor struct {
	queue taskQueue
	rate  time.Duration

	inventorySvc dgx.InventoryService
	deliverySvc  dgx.DeliveryService
}

func NewTaskProcessor(
	rate time.Duration,
	queue taskQueue,
	inventorySvc dgx.InventoryService,
	deliverySvc dgx.DeliveryService,
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

		task.Status = dgx.TaskStatusProcessing
		if err = p.queue.Update(ctx, task); err != nil {
			log.Printf("ERR! could not process task: %s", err)
			wg.Done()
			continue
		}
		log.Println("task get", task.ID, task.Type, task.Priority)

		var run func(context.Context, interface{}) error
		switch task.Type {
		case dgx.TaskTypeVerifyInventory:
			run = p.taskVerifyInventory
		case dgx.TaskTypeVerifyDelivery:
			run = p.taskVerifyDelivery
		}

		log.Println("task processing...", task.ID, task.Type)
		err = run(ctx, task.Payload)
		task.ElapsedMs = time.Since(start).Milliseconds()
		if err != nil {
			log.Printf("ERR! running tasks: %s %s", task.Type, err)
			task.Status = dgx.TaskStatusError
			task.Note = fmt.Sprintf("err: %s", err)
			if err = p.queue.Update(ctx, task); err != nil {
				log.Printf("ERR! could not run task: %s", err)
			}
			wg.Done()
			continue
		}

		task.Status = dgx.TaskStatusDone
		log.Println("task done!", task.ID, time.Duration(task.ElapsedMs)*time.Millisecond)
		if err = p.queue.Update(ctx, task); err != nil {
			log.Printf("ERR! could not update task: %s", err)
		}
		wg.Done()
	}
}

func (p *TaskProcessor) taskVerifyInventory(ctx context.Context, data interface{}) error {
	var market dgx.Market
	if err := marshallTaskPayload(data, &market); err != nil {
		return err
	}

	if market.User == nil || market.Item == nil {
		return fmt.Errorf("skipped process! missing data user:%#v item:%#v", market.User, market.Item)
	}
	if market.IsResell() {
		return nil
	}

	source := steaminvorg.InventoryAssetWithCache
	status, assets, err := verifying.Inventory(source, market.User.SteamID, market.Item.Name)
	if err != nil {
		return err
	}

	err = p.inventorySvc.Set(ctx, &dgx.Inventory{
		MarketID: market.ID,
		Status:   status,
		Assets:   assets,
	})
	return nil
}

func (p *TaskProcessor) taskVerifyDelivery(ctx context.Context, data interface{}) error {
	var market dgx.Market
	if err := marshallTaskPayload(data, &market); err != nil {
		return err
	}

	if market.User == nil || market.Item == nil {
		return fmt.Errorf("skipped process! missing data user:%#v item:%#v", market.User, market.Item)
	}

	src := steaminvorg.InventoryAsset
	status, assets, err := verifying.Delivery(src, market.User.Name, market.PartnerSteamID, market.Item.Name)
	if err != nil {
		return err
	}

	err = p.deliverySvc.Set(ctx, &dgx.Delivery{
		MarketID: market.ID,
		Status:   status,
		Assets:   assets,
	})
	return err
}

type taskQueue interface {
	Get(ctx context.Context) (*dgx.Task, error)
	Update(ctx context.Context, t dgx.Task) error
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
