package list

import (
	"fmt"
	"time"

	"github.com/dishbreak/gomilk/cli/login"

	"github.com/dishbreak/gomilk/client/task"
	"github.com/urfave/cli"

	"github.com/kyokomi/emoji"
)

type taskView struct {
	task.Task
}

func (t *taskView) String() string {
	return fmt.Sprintf("%s (due: %s) %s", t.Name(), t.dueDateString(), t.Priority())
}

func (t *taskView) dueDateString() string {

	// N/A
	// Today
	// 6:00pm
	// Tomorrow
	// Jan 3
	// Jan 3 - OVERDUE
	// Jan 4, 2020
	date, err := t.DueDate()
	if err != nil {
		return "N/A"
	}
	localZone, err := time.LoadLocation("Local")
	if err != nil {
		localZone, _ = time.LoadLocation("UTC")
	}
	now := time.Now().In(localZone)

	switch {
	case date.SameDateAs(now) && t.DueDateHasTime():
		return date.Format("3:04 PM")
	case date.SameDateAs(now):
		return "Today"
	case date.SameDateAs(now.AddDate(0, 0, 1)):
		return "Tomorrow"
	case date.Sub(now) < 0 && !t.IsCompleted() && date.SameYearAs(now):
		return date.Format("Jan 2") + emoji.Sprint(" :heavy_exclamation_mark:")
	case date.Sub(now) < 0 && !t.IsCompleted():
		return date.Format("Jan 2, 2006") + emoji.Sprint(" :heavy_exclamation_mark:")
	case date.SameYearAs(now):
		return date.Format("Jan 2")
	default:
		return date.Format("Jan 2, 2006")
	}
}

func List(c *cli.Context) error {
	args := c.Args()

	filter := "status:incomplete"
	if len(args) == 1 {
		filter = args.Get(0)
	}

	client, err := task.NewClient(login.Token)
	if err != nil {
		return err
	}

	tasks, err := client.List(filter)
	if err != nil {
		return err
	}

	for idx, task := range tasks {
		fmt.Printf("[%d] %s\n", idx, &taskView{task})
	}

	return err
}
