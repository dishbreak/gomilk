package due

import (
	"errors"
	"fmt"

	"github.com/dishbreak/gomilk/cli/login"
	"github.com/dishbreak/gomilk/cli/utils"
	"github.com/dishbreak/gomilk/client/task"
	"github.com/urfave/cli"
)

// Due handles using the setDueDate method.
func Due(c *cli.Context) error {
	args := c.Args()
	due := args[0]
	identifiers, err := utils.ResolveIdentifier(args[1:])
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
		task, err := client.SetDueDate(tasks[identifier], due)
		if err != nil {
			fmt.Printf("Error editing '%s': %s\n", tasks[identifier].Name(), err)
			errored = true
		} else {
			updatedDueDate, err := task.DueDate()
			dueDateString := "N/A"
			if err == nil {
				dueDateString = updatedDueDate.Format("Mon Jan 2")
			}
			fmt.Printf("Task '%s' is now due %s\n", task.Name(), dueDateString)
		}
	}

	if errored {
		return errors.New("Encountered errors while trying to complete tasks")
	}

	return nil
}
