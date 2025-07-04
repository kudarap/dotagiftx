package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/verify"
)

type TaskProcessor struct {
	queue taskQueue
	rate  time.Duration

	inventorySvc dotagiftx.InventoryService
	deliverySvc  dotagiftx.DeliveryService
	verify       *verify.Source
}

func NewTaskProcessor(
	rate time.Duration,
	queue taskQueue,
	inventorySvc dotagiftx.InventoryService,
	deliverySvc dotagiftx.DeliveryService,
	source *verify.Source,
) *TaskProcessor {
	return &TaskProcessor{
		queue:        queue,
		rate:         rate,
		inventorySvc: inventorySvc,
		deliverySvc:  deliverySvc,
		verify:       source,
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

		task.Status = dotagiftx.TaskStatusProcessing
		if err = p.queue.Update(ctx, task); err != nil {
			log.Printf("ERR! could not process task: %s", err)
			wg.Done()
			continue
		}
		log.Println("task get", task.ID, task.Type, task.Priority)

		var run func(context.Context, interface{}) error
		switch task.Type {
		case dotagiftx.TaskTypeVerifyInventory:
			run = p.taskVerifyInventory
		case dotagiftx.TaskTypeVerifyDelivery:
			run = p.taskVerifyDelivery
		}

		log.Println("task processing...", task.ID, task.Type)
		err = run(ctx, task.Payload)
		task.ElapsedMs = time.Since(start).Milliseconds()
		if err != nil {
			log.Printf("ERR! running tasks: %s %s", task.Type, err)
			task.Status = dotagiftx.TaskStatusError
			task.Note = fmt.Sprintf("err: %s", err)
			if err = p.queue.Update(ctx, task); err != nil {
				log.Printf("ERR! could not run task: %s", err)
			}
			wg.Done()
			continue
		}

		task.Status = dotagiftx.TaskStatusDone
		log.Println("task done!", task.ID, time.Duration(task.ElapsedMs)*time.Millisecond)
		if err = p.queue.Update(ctx, task); err != nil {
			log.Printf("ERR! could not update task: %s", err)
		}
		wg.Done()
	}
}

func (p *TaskProcessor) taskVerifyInventory(ctx context.Context, data interface{}) error {
	var market dotagiftx.Market
	if err := marshallTaskPayload(data, &market); err != nil {
		return err
	}

	if market.User == nil || market.Item == nil {
		return fmt.Errorf("skipped process! missing data user:%#v item:%#v", market.User, market.Item)
	}
	// Skips resell items.
	if market.IsResell() {
		return nil
	}

	start := time.Now()
	result, err := p.verify.Inventory(ctx, market.User.SteamID, market.Item.Name)
	if err != nil {
		return err
	}
	err = p.inventorySvc.Set(ctx, &dotagiftx.Inventory{
		MarketID:   market.ID,
		Status:     result.Status,
		Assets:     result.Assets,
		VerifiedBy: result.VerifiedBy,
		ElapsedMs:  time.Since(start).Milliseconds(),
	})
	return nil
}

func (p *TaskProcessor) taskVerifyDelivery(ctx context.Context, data interface{}) error {
	var market dotagiftx.Market
	if err := marshallTaskPayload(data, &market); err != nil {
		return err
	}

	if market.User == nil || market.Item == nil {
		return fmt.Errorf("skipped process! missing data user:%#v item:%#v", market.User, market.Item)
	}

	start := time.Now()
	result, err := p.verify.Delivery(ctx, market.User.Name, market.PartnerSteamID, market.Item.Name)
	if err != nil {
		return err
	}
	err = p.deliverySvc.Set(ctx, &dotagiftx.Delivery{
		MarketID:   market.ID,
		Status:     result.Status,
		Assets:     result.Assets,
		VerifiedBy: result.VerifiedBy,
		ElapsedMs:  time.Since(start).Milliseconds(),
	})
	return err
}

type taskQueue interface {
	Get(ctx context.Context) (*dotagiftx.Task, error)
	Update(ctx context.Context, t dotagiftx.Task) error
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
