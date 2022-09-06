package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	logging "github.com/ipfs/go-log/v2"
	"github.com/mbndr/figlet4go"

	"github.com/allisterb/citizen5/util"
)

type PingCmd struct {
	Name  string `arg:"" name:"name" help:"Object id or key name to generate 128-bit object id for."`
	Parse bool   `help:"Parse name as 128-bit object id." short:"P"`
}

var log = logging.Logger("CLI")

// Command-line arguments
var CLI struct {
	Debug bool    `help:"Enable debug mode."`
	Oid   PingCmd `cmd:"" help:"Generate or parse Motr object id."`
}

func main() {
	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		figlet4go.ColorGreen,
		figlet4go.ColorYellow,
		figlet4go.ColorCyan,
	}
	renderStr, _ := ascii.RenderOpts("citizen5", options)
	fmt.Print(renderStr)
	ctx := kong.Parse(&CLI)
	if util.Contains(ctx.Args, "--debug") {
		logging.SetAllLoggers(logging.LevelInfo)
		log.Info("Debug mode enabled.")
	}
	ctx.FatalIfErrorf(ctx.Run(&kong.Context{}))
}
