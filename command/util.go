package command

import (
	"bufio"
	"os"
	"strings"
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
