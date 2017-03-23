package command

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Dataman-Cloud/swan/src/types"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var stateMap = map[string]string{
	"slot_task_running":       "running",
	"slot_task_pending_offer": "pending",
}

var healthyMap = map[bool]string{
	true:  "true",
	false: "false",
}

// NewInspectCommand returns the CLI command for "show"
func NewInspectCommand() cli.Command {
	return cli.Command{
		Name:      "inspect",
		Usage:     "inspect application info",
		ArgsUsage: "[name]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "json",
				Usage: "List tasks with json format",
			},
			cli.BoolFlag{
				Name:  "history",
				Usage: "List task histories",
			},
		},

		Action: func(c *cli.Context) error {
			if err := inspectApplication(c); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return nil
		},
	}
}

// inspectApplication executes the "inspect" command.
func inspectApplication(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return fmt.Errorf("App ID required")
	}

	swanAddr, err := getRemote("swan")
	if err != nil {
		return err
	}

	if swanAddr == "" {
		return fmt.Errorf("swan address not found")
	}

	resp, err := http.Get(fmt.Sprintf("%s/apps/%s", swanAddr, c.Args()[0]))
	if err != nil {
		return fmt.Errorf("Unable to do request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("404")
	}

	var app *types.App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return err
	}

	data, err := json.Marshal(&app.Tasks)
	if err != nil {
		return err
	}

	if c.IsSet("json") {
		fmt.Fprintln(os.Stdout, string(data))
	} else {
		printTaskTable(app.Tasks)
	}

	return nil
}

// printTable output tasks list as table format.
func printTaskTable(tasks []*types.Task) {
	tb := tablewriter.NewWriter(os.Stdout)
	tb.SetHeader([]string{
		"Name",
		"CPUS",
		"MEM",
		"DISK",
		"IMAGE",
		"STATUS",
		"VERSIONID",
		"HISTORIES",
		"HEALTHY",
	})
	for _, task := range tasks {
		tb.Append([]string{
			task.ID,
			fmt.Sprintf("%.2f", task.CPU),
			fmt.Sprintf("%.f", task.Mem),
			fmt.Sprintf("%.f", task.Disk),
			task.Image,
			stateMap[task.Status],
			task.VersionID,
			fmt.Sprintf("%d", len(task.History)),
			healthyMap[task.Healthy],
		})
	}
	tb.Render()
}
