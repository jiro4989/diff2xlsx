package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jiro4989/diff2xlsx/command"
)

var GlobalFlags = []cli.Flag{}

var Commands = []cli.Command{
	{
		Name:   "write",
		Usage:  "Write diff to xlsx from stdin diff.",
		Action: command.CmdWrite,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "out-file-path,o",
				Value: "",
				Usage: "Output file path.",
			},
			cli.BoolFlag{
				Name:  "no-attribute,n",
				Usage: "No text attribute.",
			},
		},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
