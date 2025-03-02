package main

import (
	"github.com/alecthomas/kong"
	"github.com/jackrosenthal/nextride-go/api"
	"github.com/jackrosenthal/nextride-go/cmd"
)

func main() {
	ctx := kong.Parse(&cmd.CLI)
	context := &cmd.CliContext{
		Client: api.NewNextRideClient(),
	}
	err := ctx.Run(context)
	ctx.FatalIfErrorf(err)
}
