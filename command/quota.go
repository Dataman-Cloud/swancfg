package command

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		printSingleQuota(c, user)
	}

	return nil

}

func printSingleQuota(quota map[string]*types.Quota, user string) {
	tb := tablewriter.NewWriter(os.Stdout)
	tb.SetHeader([]string{
		"CLUSTER",
		"CPU",
		"USED",
		"MEMORY",
		"USED",
	})
	for cluster, q := range quota {
		cpuUsed, memUsed, err := getUsedQuota(user, cluster)
		if err != nil {
			fmt.Printf("calculating resource error: %s\n", err.Error())
		}
		tb.Append([]string{
			cluster,
			fmt.Sprintf("%.2f", q.Cpu),
			fmt.Sprintf("%.2f", cpuUsed),
			fmt.Sprintf("%.2f", q.Memory),
			fmt.Sprintf("%.2f", memUsed),
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
		"USED",
		"MEMORY",
		"USED",
	})
	for user, cluster := range quota {
		for c, q := range cluster {
			cpuUsed, memUsed, err := getUsedQuota(user, c)
			if err != nil {
				fmt.Printf("calculating resource error: %s\n", err.Error())
			}
			tb.Append([]string{
				user,
				c,
				fmt.Sprintf("%.2f", q.Cpu),
				fmt.Sprintf("%.2f", cpuUsed),
				fmt.Sprintf("%.2f", q.Memory),
				fmt.Sprintf("%.2f", memUsed),
			})

		}
	}
	tb.SetRowLine(true)
	tb.Render()

}

func getUsedQuota(user, cluster string) (float64, float64, error) {
	addr, err := getClusterAddr(cluster)
	if err != nil {
		return 0, 0, fmt.Errorf("Cluster can't be found. %s", err.Error())
	}

	if addr == "" {
		return 0, 0, fmt.Errorf("Cluster address can't be found. %s", cluster)
	}

	resp, err := http.Get(fmt.Sprintf("%s/apps?fields=runAs==%s", addr, user))
	if err != nil {
		return 0, 0, fmt.Errorf("Get apps failed: %s", err.Error())
	}
	defer resp.Body.Close()

	var apps []*types.App
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		return 0, 0, err
	}

	var usedCpu, usedMem float64
	for _, app := range apps {
		resp, _ := http.Get(fmt.Sprintf("%s/apps/%s", addr, app.ID))
		defer resp.Body.Close()
		var app *types.App
		if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
			return 0, 0, err
		}
		for _, task := range app.Tasks {
			usedCpu += task.Cpu
			usedMem += task.Mem
		}
	}

	return usedCpu, usedMem, nil
}
