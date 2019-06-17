package list

import (
	"fmt"

	"github.com/dishbreak/gomilk/api/tasks"
	"github.com/dishbreak/gomilk/cli/login"
	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"
)

func List(c *cli.Context) error {
	args := c.Args()

	filter := "status:incomplete"
	if len(args) == 1 {
		filter = args.Get(0)
	}

	tasks, err := tasks.GetList(login.Token, filter)
	log.WithFields(log.Fields{
		"record_count": len(tasks),
	}).Debug("Got response.")
	if err != nil {
		return err
	}

	for idx, task := range tasks {
		fmt.Printf("[%d] %s\n", idx, task.Name())
	}

	return nil
}
