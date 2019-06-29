package task

import (
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
)

// Task represents the interface that CLI will use to process API results.
type Task interface {
	// ID returns a tuple of the 3 identifiers: task, taskseries, and list.
	ID() (string, string, string)
	// Name returns the user-friendly name of the taskseries/task.
	Name() string
	// DueDate will provide the task's due date if it has one, or error if it doesn't.
	DueDate() (DateTime, error)
	// DueDateHasTime tells us if the task has a specific time it's due, or if it's just due on the day.
	DueDateHasTime() bool
	// IsCompleted will return true if the task has a completed date
	IsCompleted() bool
	// Priotity will return the priority of the task (4 is no priority)
	Priority() TaskPriority
}

/*
DateTime is a wrapper for golang's time.Time that includes some convenince methods.
*/
type DateTime struct {
	time.Time
	HasTime bool
}

/*
SameDateAs will return true if the other time.Time has the same month, day, and year.

The method will convert the other time to this DateTime location prior to making the check.
*/
func (t DateTime) SameDateAs(other time.Time) bool {
	y1, m1, d1 := t.Date()
	y2, m2, d2 := other.In(t.Location()).Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}

/*
SameDateAsDateTime will return true if the other DateTime has the same month, day, and year.

The method will convert the other time to this DateTime location prior to making the check.
*/
func (t DateTime) SameDateAsDateTime(other DateTime) bool {
	y1, m1, d1 := t.Date()
	y2, m2, d2 := other.In(t.Location()).Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}

/*
SameYearAs will return true if the other time.Time has the same Year.

The method will convert the other time to this DateTime location prior to making the check.
*/
func (t DateTime) SameYearAs(other time.Time) bool {
	return t.Year() == other.In(t.Location()).Year()
}

/*
LessThan will return true if the other time.Time describes a later time.
*/
func (t DateTime) LessThan(other DateTime) bool {
	return t.Sub(other.Time) < 0
}

/*
Sort will arrange tasks using the due date and the task name
*/
func Sort(slice []Task) {
	less := func(i, j int) (r bool) {
		defer func() {
			log.WithField("less?", r).Debug("Finished comparison.")
		}()

		one := slice[i]
		other := slice[j]

		oneDate, oneErr := one.DueDate()
		otherDate, otherErr := other.DueDate()

		switch {
		// a task with a higher priority always comes first
		case one.Priority() < other.Priority():
			r = true
		case one.Priority() > other.Priority():
			r = false
		// a task with a due date always comes first
		case oneErr == nil && otherErr != nil:
			r = true
		case oneErr != nil && otherErr == nil:
			r = false
		case oneErr != nil && otherErr != nil:
			r = one.Name() < other.Name()
		// a task with a specific time always comes first
		case oneDate.SameDateAsDateTime(otherDate) && one.DueDateHasTime() && !other.DueDateHasTime():
			r = true
		case oneDate.SameDateAsDateTime(otherDate) && !one.DueDateHasTime() && other.DueDateHasTime():
			r = false

		// tasks with identical due dates get sorted in alpha order
		case oneDate == otherDate:
			r = one.Name() < other.Name()

		default:
			r = oneDate.LessThan(otherDate)
		}
		return
	}

	log.Debug(slice)
	sort.Slice(slice, less)
}

type taskRecord struct {
	RawID        string
	RawName      string
	TaskseriesID string
	ListID       string
	Due          string
	HasDueTime   string
	Tags         []string
	Completed    string
	RawPriority  string
}

/*
Name returns the task name.
*/
func (t taskRecord) Name() string {
	return t.RawName
}

/*
DueDate returns the due date, or an error if it doesn't have one.
*/
func (t taskRecord) DueDate() (DateTime, error) {
	// Determine the timezone.
	zone, err := time.LoadLocation("Local")
	if err != nil {
		zone, err = time.LoadLocation("UTC")
		if err != nil {
			panic(err)
		}
	}

	parsed, err := time.Parse(time.RFC3339, t.Due)
	if err != nil {
		parsed = time.Now().In(zone)
	}
	return DateTime{parsed.In(zone), t.DueDateHasTime()}, err
}

/*
DueDateHasTime returns true when the user gave a specific date and time for the due date,
false when the task just has a date.
*/
func (t taskRecord) DueDateHasTime() bool {
	return t.HasDueTime == "1"
}

/*
IsCompleted returns true when the completed field is a datetime.
*/
func (t taskRecord) IsCompleted() bool {
	completed := true
	_, err := time.Parse(time.RFC3339, t.Completed)
	if err != nil {
		completed = false
	}

	return completed
}

func (t taskRecord) ID() (string, string, string) {
	return t.RawID, t.TaskseriesID, t.ListID
}
