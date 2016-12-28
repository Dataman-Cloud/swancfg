package command

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Dataman-Cloud/swancfg/types"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func NewListCommand() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list all applications",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "json",
				Usage: "List apps with json format",
			},
			cli.BoolFlag{
				Name:  "all",
				Usage: "List all apps",
			},
			cli.StringFlag{
				Name:  "cluster",
				Usage: "List apps belong to cluster [CLUSTER]",
			},
			cli.StringFlag{
				Name:  "user",
				Usage: "List apps belong to user [USER]",
			},
		},
		Action: func(c *cli.Context) error {
			if err := listApps(c); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return nil
		},
	}

}

func listApps(c *cli.Context) error {

	var apps []*types.App

	cluster, user := c.String("cluster"), c.String("user")

	if cluster == "" && user == "" {
		apps, _ = getAllApps("")
	}

	if cluster != "" && user == "" {
		apps, _ = getAppsByClusterID(cluster)
	}

	if cluster == "" && user != "" {
		apps, _ = getAppsByUser(user)
	}

	if cluster != "" && user != "" {
		apps, _ = getAppsByClusterAndUser(cluster, user)
	}

	printTable(apps)

	return nil
}

func getClusters() ([]string, error) {
	f, err := os.Open("cluster.cfg")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	clusters := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "\t\t")
		clusters = append(clusters, line[1])
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return clusters, nil
}

func getCluster(cluster string) (string, error) {
	f, err := os.Open("cluster.cfg")
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "\t\t")
		if line[0] == cluster {
			return line[1], nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func getAllApps(filter string) ([]*types.App, error) {
	clusters, err := getClusters()
	if err != nil {
		return nil, err
	}

	var allApps []*types.App
	for _, cluster := range clusters {
		resp, _ := http.Get(fmt.Sprintf("%s/apps/?fields=%s", cluster, filter))
		defer resp.Body.Close()
		var apps []*types.App
		if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
			continue
		}

		allApps = append(allApps, apps...)
	}

	return allApps, nil
}

func getAppsByClusterID(clusterId string) ([]*types.App, error) {
	cluster, err := getCluster(clusterId)
	if err != nil {
		return nil, err
	}

	if cluster == "" {
		return nil, fmt.Errorf("cluster not found")
	}

	resp, _ := http.Get(fmt.Sprintf("%s/apps/", cluster))
	defer resp.Body.Close()
	var apps []*types.App
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		return nil, err
	}

	return apps, nil
}

func getAppsByUser(user string) ([]*types.App, error) {
	fmt.Println("get apps by user")
	filter := fmt.Sprintf("runAs==%s", user)

	return getAllApps(filter)
}

func getAppsByClusterAndUser(clusterId, userId string) ([]*types.App, error) {
	apps, _ := getAppsByClusterID(clusterId)

	var results []*types.App
	for _, app := range apps {
		if app.RunAs == userId {
			results = append(results, app)
		}
	}

	return results, nil
}

func printTable(apps []*types.App) {
	tb := tablewriter.NewWriter(os.Stdout)
	tb.SetHeader([]string{
		"ID",
		"Name",
		"Instances",
		"RunAS",
		"ClusterID",
		"State",
		"Created",
		"Updated",
	})
	for _, app := range apps {
		tb.Append([]string{
			app.ID,
			app.Name,
			fmt.Sprintf("%d", app.Instances),
			app.RunAs,
			app.ClusterId,
			app.State,
			app.Created.Format("2006-01-02 15:04:05"),
			app.Updated.Format("2006-01-02 15:04:05"),
		})
	}
	tb.Render()
}
