package command

import (
	"fmt"
	"io/ioutil"

	"github.com/Dataman-Cloud/swancfg/types"
	"gopkg.in/yaml.v2"
)

func getRemotes() (map[string]string, error) {
	db, err := NewBoltStore(".bolt.db")
	if err != nil {
		return nil, fmt.Errorf("Init store engine failed:%s", err)
	}

	tx, err := db.conn.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte("swan"))

	c := bucket.Cursor()

	remotes := make(map[string]string)
	for k, v := c.First(); k != nil; k, v = c.Next() {
		remotes[string(k)] = string(v)
	}

	return remotes, nil
}

func getRemote(remote string) (string, error) {
	db, err := NewBoltStore(".bolt.db")
	if err != nil {
		return "", fmt.Errorf("Init store engine failed:%s", err)
	}

	tx, err := db.conn.Begin(false)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte("swan"))
	val := bucket.Get([]byte(remote))

	if val == nil {
		return "", nil
	}

	return string(val[:]), nil
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
