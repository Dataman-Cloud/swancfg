package command

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func NewClusterCommand() cli.Command {
	return cli.Command{
		Name:  "cluster",
		Usage: "cluster manage",
		Subcommands: []cli.Command{
			cli.Command{
				Name:  "list",
				Usage: "list clusters",
				Action: func(c *cli.Context) {
					if err := listClusters(c); err != nil {
						fmt.Fprintln(os.Stderr, "Error:", err)
					}
				},
			},
			cli.Command{
				Name:      "add",
				Usage:     "add cluster",
				ArgsUsage: "[name] [address]",
				Action: func(c *cli.Context) {
					if err := addCluster(c); err != nil {
						fmt.Fprintln(os.Stderr, "Error:", err)
					}
				},
			},
		},
	}
}

func listClusters(c *cli.Context) error {
	clusters, err := getClusters()
	if err != nil {
		return err
	}

	printClusters(clusters)

	return nil
}

func printClusters(clusters map[string]string) {
	tb := tablewriter.NewWriter(os.Stdout)
	tb.SetHeader([]string{
		"CLUSTER",
		"ADDRESS",
	})
	for name, addr := range clusters {
		tb.Append([]string{
			name,
			addr,
		})
	}
	tb.SetRowLine(true)
	tb.Render()
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
