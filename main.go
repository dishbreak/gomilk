package main

import (
	"fmt"
	"os"

	"github.com/dishbreak/gomilk/cli/login"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gomilk"
	app.Usage = "Remember the Milk command-line client."
	app.Commands = []cli.Command{
		{
			Name:    "login",
			Aliases: []string{"l"},
			Usage:   "Log in to Remember the Milk.",
			Action:  login.Login,
		},
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
