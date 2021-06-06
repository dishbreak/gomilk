package main

import (
	"fmt"
	"os"

	"github.com/dishbreak/gomilk/api"
	"github.com/dishbreak/gomilk/cli/add"
	"github.com/dishbreak/gomilk/cli/complete"
	"github.com/dishbreak/gomilk/cli/due"
	"github.com/dishbreak/gomilk/cli/list"
	"github.com/dishbreak/gomilk/cli/login"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	api.ValidateCredentials()

	app := cli.NewApp()
	app.Name = "gomilk"
	app.Usage = "Remember the Milk command-line client."
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Enable verbose output for debugging",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "login",
			Aliases: []string{"l"},
			Usage:   "Log in to Remember the Milk.",
			Action:  login.Login,
		},
		{
			Name:    "add",
			Aliases: []string{"t"},
			Usage:   "Add a task to Remember the Milk.",
			Action:  makeAction(add.Add),
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list tasks using an optional filter.",
			Action:  makeAction(list.List),
		},
		{
			Name:    "complete",
			Aliases: []string{"c"},
			Usage:   "completes listed tasks",
			Action:  makeAction(complete.Complete),
		},
		{
			Name:      "due",
			Aliases:   []string{"d"},
			Usage:     "sets due date for listed tasks",
			ArgsUsage: "DUE_DATE identifiers...",
			Action:    makeAction(due.Due),
		},
	}

	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			fmt.Println("turning on debug output")
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func makeAction(f func(*cli.Context) error) func(*cli.Context) error {
	return func(c *cli.Context) error {
		if !login.IsAuthenticated() {
			return fmt.Errorf("either you're not logged in or your credentials are expired. Run `%s login` first, then try again", c.App.Name)
		}
		err := login.Setup()
		if err != nil {
			return err
		}
		return f(c)
	}
}
