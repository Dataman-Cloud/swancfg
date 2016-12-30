package command

import (
	"fmt"
	"os"

	"github.com/Dataman-Cloud/swancfg/types"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func NewQuotaCommand() cli.Command {
	return cli.Command{
		Name:  "quota",
		Usage: "quota manage",
		Subcommands: []cli.Command{
			cli.Command{
				Name:  "list",
				Usage: "show quota",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "user",
						Usage: "list quota for user [USER]",
					},
				},
				Action: func(c *cli.Context) {
					if err := listQuota(c); err != nil {
						fmt.Fprintln(os.Stderr, "Error:", err)
					}
				},
			},
		},
	}
}

func listQuota(c *cli.Context) error {
	if c.String("user") != "" {
		return listQuotaForUser(c.String("user"))
	}

	return listQuotaForAll()
}

func listQuotaForUser(user string) error {
	quota, err := getQuotas()
	if err != nil {
		return err
	}

	if c, ok := quota[user]; ok {
		printSingleQuota(c)
	}

	return nil

}

func printSingleQuota(quota map[string]*types.Quota) {
	tb := tablewriter.NewWriter(os.Stdout)
	tb.SetHeader([]string{
		"CLUSTER",
		"CPU",
		"MEMORY",
	})
	for cluster, q := range quota {
		tb.Append([]string{
			cluster,
			fmt.Sprintf("%.2f", q.Cpu),
			fmt.Sprintf("%.2f", q.Memory),
		})
	}
	tb.SetRowLine(true)
	tb.Render()
}

func listQuotaForAll() error {
	quota, err := getQuotas()
	if err != nil {
		return err
	}

	printAllQuota(quota)

	return nil
}

func printAllQuota(quota Quota) {
	tb := tablewriter.NewWriter(os.Stdout)
	tb.SetHeader([]string{
		"USER",
		"CLUSTER",
		"CPU",
		"MEMORY",
	})
	for user, cluster := range quota {
		for c, q := range cluster {
			tb.Append([]string{
				user,
				c,
				fmt.Sprintf("%.2f", q.Cpu),
				fmt.Sprintf("%.2f", q.Memory),
			})

		}
	}
	tb.SetRowLine(true)
	tb.Render()

}
