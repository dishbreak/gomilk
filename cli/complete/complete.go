package complete

import (
	"github.com/dishbreak/gomilk/cli/utils"
	"github.com/urfave/cli"
)



func Complete(c *cli.Context) error {
	args := c.Args()

	identifiers, err := utils.ResolveIdentifier(args)
	if err != nil {
		return err
	}

	cache, err := utils.NewCache("task")
	if err != nil {
		return err
	}


	records, err := cache.Get()
	if err != nil {
		return err
	}

	if len

	for _, identifier 

}
