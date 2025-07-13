package models

import "time"

// Representing task table in database
type Task struct {
	Id        int64
	Status    TaskStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TaskStatus string

const (
	New        TaskStatus = "NEW"        // Default status. task not finished yet and not being porocessed right now niether.
	Processing TaskStatus = "PROCESSING" // Task is being processed.
	Processed  TaskStatus = "PROCESSED"  // Task is in final status.
)
