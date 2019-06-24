package add

import (
	"fmt"

	"github.com/dishbreak/gomilk/cli/utils"

	"github.com/dishbreak/gomilk/api/tasks"
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

	timeline, err := utils.Timeline()
	if err != nil {
		return err
	}

	task, err := tasks.Add(login.Token, args.Get(0), timeline)
	if err != nil {
		return err
	}

	if dueDate, err := task.DueDate(); err != nil {
		fmt.Printf("Created task '%s', due on '%s'\n", task.Name(), dueDate)
	} else {
		fmt.Printf("Created task '%s'\n", task.Name())
	}

	return nil
}
