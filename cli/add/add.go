package add

import (
	"fmt"

	"github.com/dishbreak/gomilk/client/task"

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

	client, err := task.NewClient(login.Token)
	if err != nil {
		return err
	}

	taskObj, err := client.Add(args[0])
	if err != nil {
		return err
	}

	if dueDate, err := taskObj.DueDate(); err != nil {
		fmt.Printf("Created task '%s', due on '%s'\n", taskObj.Name(), dueDate)
	} else {
		fmt.Printf("Created task '%s'\n", taskObj.Name())
	}

	return nil
}
