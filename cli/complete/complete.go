package complete

import (
	"errors"
	"fmt"

	"github.com/dishbreak/gomilk/cli/login"
	"github.com/dishbreak/gomilk/cli/utils"
	"github.com/dishbreak/gomilk/client/task"
	"github.com/urfave/cli"
)

/*
Complete will attempt to mark given tasks complete
*/
func Complete(c *cli.Context) error {
	args := c.Args()

	identifiers, err := utils.ResolveIdentifier(args)
	if err != nil {
		return err
	}

	client, err := task.NewClient(login.Token)
	if err != nil {
		return err
	}

	tasks, err := client.GetCachedTasks()

	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		return errors.New("Use the list command to display tasks first")
	}
	if maxIdentifier, maxTask := identifiers[len(identifiers)-1], len(tasks)-1; maxIdentifier > maxTask {
		return fmt.Errorf("Identifer %d is greater than highest task number (%d)", maxIdentifier, maxTask)
	}

	errored := false
	for _, identifier := range identifiers {
		task, err := client.Complete(tasks[identifier])
		if err != nil {
			fmt.Printf("Error completing '%s': %s\n", tasks[identifier].Name(), err)
			errored = true
		} else {
			fmt.Printf("Completed task '%s'\n", task.Name())
		}
	}

	if errored {
		return errors.New("Encountered errors while trying to complete tasks")
	}

	return nil
}
