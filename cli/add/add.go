package add

import (
	"fmt"

	"github.com/dishbreak/gomilk/api/tasks"
	"github.com/dishbreak/gomilk/api/timelines"
	"github.com/dishbreak/gomilk/cli/login"

	"github.com/urfave/cli"
)

/*
Add uses the RTM API to create a new task.
*/
func Add(c *cli.Context) error {
	args := c.Args()

	if len(args) != 1 {
		return fmt.Errorf("need exactly one argument (got %d)", len(args))
	}

	timeline, err := timelines.Create(login.Token)
	if err != nil {
		return err
	}

	task, err := tasks.Add(login.Token, args.Get(0), timeline.Rsp.Timeline)
	if err != nil {
		return err
	}

	if task.Due() != nil {
		fmt.Printf("Created task '%s', due on '%s'\n", task.Name(), task.Due())
	} else {
		fmt.Printf("Created task '%s'\n", task.Name())
	}

	return nil
}
