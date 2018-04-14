package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/jiro4989/diff2xlsx/internal/version"
)

func main() {
	app := cli.NewApp()
	app.Name = version.Name
	app.Version = version.Version + " " + version.Revision
	app.Author = "jiro4989"
	app.Email = ""
	app.Usage = "Convert diff to xlsx."

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}
