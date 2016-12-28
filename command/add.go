package command

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

func NewAddCommand() cli.Command {
	return cli.Command{
		Name:      "add-cluster",
		Usage:     "add swan cluster",
		ArgsUsage: "[name] [address]",
		Action: func(c *cli.Context) error {
			if err := addCluster(c); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return nil
		},
	}
}

func addCluster(c *cli.Context) error {
	if len(c.Args()) < 2 {
		fmt.Println("Missing argument")
		fmt.Println("swancfg add-cluster [name] [address]")
		return nil
	}

	path := "./cluster.cfg"

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("Open cluster.cfg failed: %s", err.Error())
	}

	content := fmt.Sprintf("%s\t\t%s\n", c.Args()[0], c.Args()[1])
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("Write to file cluster.cfg failed: %s", err.Error())
	}

	if err := file.Sync(); err != nil {
		return err
	}

	fmt.Println("add cluster successful")

	return nil
}
