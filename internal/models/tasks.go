package models

import "time"

// Representing task table in database
type Task struct {
	Id        int64      `json:"id"`
	Status    TaskStatus `json:"status"`
	UpdatedAt time.Time
	CreatedAt time.Time
}

type TaskStatus string

func (ts TaskStatus) IsValid() bool {
	switch ts {
	case New, Processing, Processed:
		return true
	default:
		return false
	}
}

const (
	New        TaskStatus = "NEW"        // Default status. task not finished yet and not being porocessed right now niether.
	Processing TaskStatus = "PROCESSING" // Task is being processed.
	Processed  TaskStatus = "PROCESSED"  // Task is in final status.
)
