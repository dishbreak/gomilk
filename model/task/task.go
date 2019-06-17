package task

import "time"

// Task represents the interface that CLI will use to process API results.
type Task interface {
	// Name provides the user-visible name of the task. This may not be unique.
	Name() string
	// DueDate will provide the task's due date if it has one, or error if it doesn't.
	DueDate() (time.Time, error)
	// DueDateHasTime tells us if the task has a specific time it's due, or if it's just due on the day.
	DueDateHasTime() bool
	// IsCompleted will return true if the task has a completed date
	IsCompleted() bool
}
