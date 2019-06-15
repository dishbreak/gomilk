package main

import (
	"fmt"
	"os"

	"github.com/dishbreak/gomilk/cli/add"
	"github.com/dishbreak/gomilk/cli/login"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
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
			Aliases: []string{"a"},
			Usage:   "Add a task to Remember the Milk.",
			Action:  makeAction(add.Add),
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
