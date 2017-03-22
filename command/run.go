package command

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Dataman-Cloud/swancfg/types"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

var (
	clusterAddr string
)

func NewRunCommand() cli.Command {
	return cli.Command{
		Name:  "run",
		Usage: "run new application",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "from-file, f",
				Usage: "Run application from `FILE`",
			},
			cli.StringFlag{
				Name:  "name, n",
				Usage: "Set application name",
			},
			cli.IntFlag{
				Name:  "times",
				Usage: "Concurrent for testing",
			},
			cli.BoolTFlag{
				Name:  "disable-quota",
				Usage: "Disable quota check",
			},
		},
		Action: func(c *cli.Context) error {
			if err := runApplication(c); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return nil
		},
	}
}

func runApplication(c *cli.Context) error {
	if c.String("from-file") == "" {
		return fmt.Errorf("Spec file must be specified for running application")
	}

	var spec *types.Spec

	file, err := ioutil.ReadFile(c.String("from-file"))
	if err != nil {
		return fmt.Errorf("Read json file failed: %s", err.Error())
	}

	if err := json.Unmarshal(file, &spec); err != nil {
		return fmt.Errorf("Unmarshal error: %s", err.Error())
	}

	name := c.String("name")
	if name != "" {
		spec.AppName = name
	}

	if !c.BoolT("disable-quota") {
		if err := checkQuota(spec); err != nil {
			return err
		}
	}

	fmt.Printf("===> sending request to cluster:%s...", spec.Cluster)
	if err := sendRequest(spec); err != nil {
		return err
	}
	fmt.Println("done")

	fmt.Printf("===> waiting for application %s-%s-%s to running...", name, spec.RunAs, spec.Cluster)
	ticker := time.NewTicker(time.Duration(1 * time.Second))
	timeout := time.NewTicker(time.Duration(10 * time.Second))
	for {
		select {
		case <-ticker.C:
			fmt.Printf(".")
			status, err := getStatus(fmt.Sprintf("%s-%s-%s", spec.AppName, spec.RunAs, spec.Cluster))
			if err != nil {
				return err
			}

			if status == "normal" {
				fmt.Printf("running\n")
				return nil
			}
		case <-timeout.C:
			fmt.Printf("timeout\n")
			return nil
		}
	}

	return nil
}

func checkQuota(spec *types.Spec) error {
	addr, err := getClusterAddr(spec.Cluster)
	if err != nil {
		return fmt.Errorf("Cluster can't be found. %s", err.Error())
	}

	if addr == "" {
		return fmt.Errorf("Cluster address can't be found. %s", spec.Cluster)
	}

	clusterAddr = addr

	resp, err := http.Get(fmt.Sprintf("%s/apps?fields=runAs==%s", clusterAddr, spec.RunAs))
	if err != nil {
		return fmt.Errorf("Get apps failed: %s", err.Error())
	}
	defer resp.Body.Close()

	var apps []*types.App
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		return err
	}

	fmt.Printf("===> calculating total used resources...\n")
	var usedCpu, usedMem float64
	for _, app := range apps {
		resp, _ := http.Get(fmt.Sprintf("%s/apps/%s", clusterAddr, app.ID))
		defer resp.Body.Close()
		var app *types.App
		if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
			return err
		}
		for _, task := range app.Tasks {
			usedCpu += task.Cpu
			usedMem += task.Mem
		}
	}

	fmt.Printf("===> calculating quota...\n")
	quota, err := getQuota(spec.RunAs, spec.Cluster)
	if err != nil {
		return fmt.Errorf("calculate quota got error: %s", err.Error())
	}

	if quota == nil {
		return fmt.Errorf("No quota found")
	}

	needCpu := float64(spec.Instances) * spec.Cpus
	needMem := float64(spec.Instances) * spec.Mem

	if (quota.Cpu-usedCpu) < needCpu || (quota.Memory-usedMem) < needMem {
		fmt.Printf("===> quota exceed...\n")
		fmt.Printf("  Total quota == Cpu: %.2f Memory: %.2f\n", quota.Cpu, quota.Memory)
		fmt.Printf("  Use quota == Cpu: %.2f Memory: %.2f\n", usedCpu, usedMem)
		fmt.Printf("  Left quota == Cpu: %.2f Memory: %.2f\n", quota.Cpu-usedCpu, quota.Memory-usedMem)
		fmt.Printf("  Need quota == Cpu: %.2f Memory: %.2f\n", needCpu, needMem)

		return fmt.Errorf("Quota exceed")
	}

	fmt.Printf("===> quota satisfied...\n")
	return nil
}

func getQuota(user, cluster string) (*types.Quota, error) {
	file, err := ioutil.ReadFile("./quota.yml")
	if err != nil {
		return nil, fmt.Errorf("Read quota file failed: %s", err.Error())
	}

	entries := make(map[string]map[string]*types.Quota)

	yaml.Unmarshal(file, entries)

	if c, ok := entries[user]; ok {
		if quota, ok := c[cluster]; ok {
			return quota, nil
		}
	}

	return nil, nil
}

func sendRequest(spec *types.Spec) error {
	addr, err := getClusterAddr(spec.Cluster)
	if err != nil {
		return fmt.Errorf("Cluster can't be found. %s", err.Error())
	}

	if addr == "" {
		return fmt.Errorf("Cluster address can't be found. %s", spec.Cluster)
	}

	clusterAddr = addr

	payload, err := json.Marshal(&spec)
	if err != nil {
		return fmt.Errorf("Marsh failed: %s", err.Error())
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/apps", clusterAddr), bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("Make new request failed: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "swancfg/0.1")

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("Send post request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Send request ok but status code not 201")
	}

	return nil
}

func getStatus(appId string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/apps/%s", clusterAddr, appId))
	if err != nil {
		return "unknown", err
	}
	defer resp.Body.Close()
	var app *types.App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return "", err
	}

	return app.State, nil
}

func getClusterAddr(name string) (string, error) {
	f, err := os.Open("cluster.cfg")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "\t\t")
		if line[0] == name {
			return line[1], nil
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return "", nil
}
