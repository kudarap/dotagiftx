package dgx

import "time"

// Task kinds.
const (
	TaskTypeVerifyDelivery  TaskType = 1
	TaskTypeVerifyInventory TaskType = 2
)

// Task priorities.
const (
	TaskPriorityHigh   TaskPriority = 1
	TaskPriorityMedium TaskPriority = 2
	TaskPriorityLow    TaskPriority = 3
)

// Task status.
const (
	TaskStatusPending    TaskStatus = 0
	TaskStatusProcessing TaskStatus = 1
	TaskStatusDone       TaskStatus = 2
	TaskStatusError      TaskStatus = 6
)

type (
	// TaskType represents task kind.
	TaskType uint8

	// TaskStatus represent task status.
	TaskStatus uint8

	// TaskPriority represent task priority.
	TaskPriority uint8

	// Task represents task data model.
	Task struct {
		ID        string       `json:"id"           db:"id,omitempty,index"`
		Status    TaskStatus   `json:"status"       db:"status,index"`
		Priority  TaskPriority `json:"priority"     db:"priority,omitempty,index"`
		Type      TaskType     `json:"type"         db:"type,omitempty,index"`
		Payload   interface{}  `json:"payload"      db:"payload,omitempty"`
		Retry     int          `json:"retry"        db:"retry,omitempty"`
		Note      string       `json:"note"         db:"note,omitempty"`
		ElapsedMs int64        `json:"elapsed_ms"   db:"elapsed_ms,omitempty"`
		CreatedAt *time.Time   `json:"created_at"   db:"created_at,omitempty,index"`
		UpdatedAt *time.Time   `json:"updated_at"   db:"updated_at,omitempty"`
	}
)

var taskTypeStrings = map[TaskType]string{
	TaskTypeVerifyDelivery:  "verify_delivery",
	TaskTypeVerifyInventory: "verify_inventory",
}

var taskPriorityStrings = map[TaskPriority]string{
	TaskPriorityHigh:   "high",
	TaskPriorityMedium: "medium",
	TaskPriorityLow:    "low",
}

var TaskStatusStrings = map[TaskStatus]string{
	TaskStatusPending:    "pending",
	TaskStatusProcessing: "processing",
	TaskStatusDone:       "done",
	TaskStatusError:      "error",
}

func (t TaskType) String() string { return taskTypeStrings[t] }

func (p TaskPriority) String() string { return taskPriorityStrings[p] }

func (s TaskStatus) String() string { return TaskStatusStrings[s] }
