package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
)

var (
	version = "dev"
	// commit  = "none"
	// date    = "unknown"
)

var cli struct {
	BuildCmd BuildCmd `cmd:"" name:"build" help:"Build the plugin JAR file."`
	YmlCmd   YmlCmd   `cmd:"" name:"yml" help:"Generate the plugin.yml file."`
}

func rootContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	return signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
