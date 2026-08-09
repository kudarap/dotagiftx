package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	dotagiftx2 "github.com/kudarap/dotagiftx/dotagiftx"
	"github.com/kudarap/dotagiftx/steam"
	"github.com/kudarap/dotagiftx/verify"
)

type TaskProcessor struct {
	queue taskQueue
	rate  time.Duration

	inventorySvc         inventoryService
	deliverySvc          deliveryService
	verify               *verify.Source
	inventoryInvalidator inventoryInvalidator

	logger *slog.Logger
}

func NewTaskProcessor(
	rate time.Duration,
	queue taskQueue,
	inventorySvc inventoryService,
	deliverySvc deliveryService,
	source *verify.Source,
	invInvalidator inventoryInvalidator,
	logger *slog.Logger,
) *TaskProcessor {
	return &TaskProcessor{
		queue:                queue,
		rate:                 rate,
		inventorySvc:         inventorySvc,
		deliverySvc:          deliverySvc,
		verify:               source,
		inventoryInvalidator: invInvalidator,
		logger:               logger.With("processor", "task"),
	}
}

func (p *TaskProcessor) Run(wg *sync.WaitGroup) {
	ctx := context.Background()
	for {
		time.Sleep(p.rate)

		start := time.Now()

		t, err := p.queue.Get(ctx)
		if err != nil {
			p.logger.ErrorContext(ctx, "get task from queue", "err", err)
			continue
		}
		if t == nil {
			continue
		}

		task := *t
		wg.Add(1)

		task.Status = dotagiftx2.TaskStatusProcessing
		if err = p.queue.Update(ctx, task); err != nil {
			p.logger.ErrorContext(ctx, "mark task status as processing", "err", err)
			wg.Done()
			continue
		}
		p.logger.InfoContext(ctx, "processing",
			"id", task.ID,
			"type", task.Type,
			"priority", task.Priority,
		)

		var run func(context.Context, any) error
		switch task.Type {
		case dotagiftx2.TaskTypeVerifyInventory:
			run = p.taskVerifyInventory
		case dotagiftx2.TaskTypeVerifyDelivery:
			run = p.taskVerifyDelivery
		}

		err = run(ctx, task.Payload)
		task.ElapsedMs = time.Since(start).Milliseconds()
		if err != nil {
			p.logger.ErrorContext(ctx, "run task",
				"id", task.ID,
				"err", err,
			)
			task.Status = dotagiftx2.TaskStatusError
			task.Note = fmt.Sprintf("err: %s", err)
			if err = p.queue.Update(ctx, task); err != nil {
				p.logger.ErrorContext(ctx, "mark task status as error",
					"id", task.ID,
					"err", err,
				)
			}
			wg.Done()
			continue
		}

		task.Status = dotagiftx2.TaskStatusDone
		p.logger.InfoContext(ctx, "done",
			"id", task.ID,
			"type", task.Type,
			"elapsed_ms", task.ElapsedMs,
		)
		if err = p.queue.Update(ctx, task); err != nil {
			p.logger.ErrorContext(ctx, "mark task status as done",
				"id", task.ID,
				"err", err,
			)
		}
		wg.Done()
	}
}

func (p *TaskProcessor) taskVerifyInventory(ctx context.Context, data any) error {
	var market dotagiftx2.Market
	if err := marshallTaskPayload(data, &market); err != nil {
		return err
	}
	if market.User == nil || market.Item == nil {
		return fmt.Errorf("skipped process! missing data user:%#v item:%#v", market.User, market.Item)
	}
	// Skips resold items.
	if market.IsResell() {
		return nil
	}

	start := time.Now()
	result, err := p.verify.Inventory(ctx, market.User.SteamID, market.Item.Name)
	if err != nil {
		return err
	}
	return p.inventorySvc.Set(ctx, &dotagiftx2.Inventory{
		MarketID:   market.ID,
		Status:     result.Status,
		Assets:     result.Assets,
		VerifiedBy: result.VerifiedBy,
		ElapsedMs:  time.Since(start).Milliseconds(),
	})
}

func (p *TaskProcessor) taskVerifyDelivery(ctx context.Context, data any) error {
	var market dotagiftx2.Market
	if err := marshallTaskPayload(data, &market); err != nil {
		return err
	}
	if err := steam.ValidateSteamID(market.PartnerSteamID); err != nil {
		return err
	}
	if market.User == nil || market.Item == nil {
		return fmt.Errorf("skipped process! missing data user:%#v item:%#v", market.User, market.Item)
	}
	if err := p.inventoryInvalidator.Invalidate(ctx, market.PartnerSteamID); err != nil {
		return fmt.Errorf("invalidate inventory: %s", err)
	}

	start := time.Now()
	result, err := p.verify.Delivery(ctx, market.User.Name, market.PartnerSteamID, market.Item.Name)
	if err != nil {
		return err
	}
	err = p.deliverySvc.Set(ctx, &dotagiftx2.Delivery{
		MarketID:   market.ID,
		Status:     result.Status,
		Assets:     result.Assets,
		VerifiedBy: result.VerifiedBy,
		ElapsedMs:  time.Since(start).Milliseconds(),
	})
	return err
}

type taskQueue interface {
	Get(ctx context.Context) (*dotagiftx2.Task, error)
	Update(ctx context.Context, t dotagiftx2.Task) error
}

type inventoryInvalidator interface {
	Invalidate(ctx context.Context, steamID string) error
}

// deliveryService provides access to delivery service methods used by the task processor.
type deliveryService interface {
	// Set saves new Delivery details.
	Set(context.Context, *dotagiftx2.Delivery) error
}

// inventoryService provides access to inventory service methods used by the task processor.
type inventoryService interface {
	// Set saves new Inventory details.
	Set(context.Context, *dotagiftx2.Inventory) error
}

func marshallTaskPayload(in, out any) error {
	raw, ok := in.(map[string]any)
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
