package command

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func NewRemoteCommand() cli.Command {
	return cli.Command{
		Name:  "remote",
		Usage: "remote address management",
		Subcommands: []cli.Command{
			cli.Command{
				Name:  "list",
				Usage: "list remote address(es)",
				Action: func(c *cli.Context) {
					if err := listRemotes(c); err != nil {
						fmt.Fprintln(os.Stderr, "Error:", err)
					}
				},
			},
			cli.Command{
				Name:      "add",
				Usage:     "add remote",
				ArgsUsage: "[swan|mesos] [address]",
				Action: func(c *cli.Context) {
					if err := addRemote(c); err != nil {
						fmt.Fprintln(os.Stderr, "Error:", err)
					}
				},
			},
		},
	}
}

func listRemotes(c *cli.Context) error {
	remotes, err := getRemotes()
	if err != nil {
		return err
	}

	printRemotes(remotes)

	return nil
}

func printRemotes(clusters map[string]string) {
	tb := tablewriter.NewWriter(os.Stdout)
	tb.SetHeader([]string{
		"REMOTE",
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

func addRemote(c *cli.Context) error {
	if len(c.Args()) < 2 {
		fmt.Println("Missing argument")
		fmt.Println("swancfg remote add [swan|mesos] [address]")
		return nil
	}

	if c.Args()[0] != "swan" && c.Args()[0] != "mesos" {
		fmt.Println("Only swan| mesos are supported")
	}

	db, err := NewBoltStore(".bolt.db")
	if err != nil {
		return fmt.Errorf("Init store engine failed:%s", err)
	}

	tx, err := db.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte("swan"))
	if err := bucket.Put([]byte(c.Args()[0]), []byte(c.Args()[1])); err != nil {
		return err
	}

	return tx.Commit()
}
