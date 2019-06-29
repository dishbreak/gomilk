package task

import (
	"strconv"

	"github.com/kyokomi/emoji"
)

type TaskPriority int

const (
	HighPriority TaskPriority = 1
	MedPriority  TaskPriority = 2
	LowPriority  TaskPriority = 3
	NoPriority   TaskPriority = 4
)

func (p TaskPriority) String() (v string) {
	switch p {
	case HighPriority:
		v = emoji.Sprint(":one:")
	case MedPriority:
		v = emoji.Sprint(":two:")
	case LowPriority:
		v = emoji.Sprint(":three:")
	case NoPriority:
		v = ""
	}
	return
}

func (t taskRecord) Priority() TaskPriority {
	priority, err := strconv.Atoi(t.RawPriority)
	if err != nil {
		return NoPriority
	}

	return TaskPriority(priority)
}
