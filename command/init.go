package command

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

func NewInitCommand() cli.Command {
	return cli.Command{
		Name:  "init",
		Usage: "init database",
		Action: func(c *cli.Context) error {
			if err := initDB(c); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			return nil
		},
	}
}

func initDB(c *cli.Context) error {
	return nil
}
