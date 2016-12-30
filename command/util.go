package command

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Dataman-Cloud/swancfg/types"
	"gopkg.in/yaml.v2"
)

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

type Quota map[string]map[string]*types.Quota

func getQuotas() (Quota, error) {
	file, err := ioutil.ReadFile("./quota.yml")
	if err != nil {
		return nil, fmt.Errorf("Read quota file failed: %s", err.Error())
	}
	entries := make(map[string]map[string]*types.Quota)

	if err := yaml.Unmarshal(file, entries); err != nil {
		return nil, fmt.Errorf("Unmarshal failed: %s", err.Error())
	}

	return entries, nil
}
