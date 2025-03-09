package main

import (
	"github.com/alecthomas/kong"
	"github.com/jackrosenthal/gtfs-cli/cmd"
)

func main() {
	ctx := kong.Parse(&cmd.CLI)
	context := &cmd.CliContext{}
	err := ctx.Run(context)
	ctx.FatalIfErrorf(err)
}
