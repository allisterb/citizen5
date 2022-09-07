package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	logging "github.com/ipfs/go-log/v2"
	"github.com/jeandeaual/go-locale"
	"github.com/mbndr/figlet4go"

	"github.com/allisterb/citizen5/nym"
	"github.com/allisterb/citizen5/util"
)

type PingCmd struct {
	Address string `arg:"" name:"address" help:"Nym address to send test message to." default:""`
	Binary  bool   `help:"Send a binary file as the test message."`
}

var log = logging.Logger("main")

// Command-line arguments
var CLI struct {
	Debug bool    `help:"Enable debug mode."`
	WSUrl string  `help:"The URL of the Nym websocket client." default:"ws://localhost:1977"`
	Ping  PingCmd `cmd:"" help:"Send a test message to a Nym address."`
}

func init() {
	logging.SetAllLoggers(logging.LevelInfo)
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

func (l *PingCmd) Run(ctx *kong.Context) error {
	lc, _ := locale.GetLanguage()
	log.Infof("Locale %s", lc)
	conn, err := nym.GetConn(CLI.WSUrl)
	if err != nil {
		log.Errorf("could not open WebSocket connection to %s:%v", CLI.WSUrl, err)
		return nil
	}
	defer conn.Close()
	if l.Address == "" {
		l.Address = nym.GetSelfAddressText(conn)
		d := len(l.Address)
		log.Info("pinging own client address...", d)
	}
	if l.Binary {
		if err := nym.SendBinary(conn, nym.GetSelfAddressBinary(conn), "main.go"); err != nil {
			log.Errorf("could not send binary message to Nym address %s:%v", l.Address, err)
			return err
		}
	} else {
		if err := nym.SendText(conn, l.Address, "hello", true); err != nil {
			log.Errorf("could not send text message to Nym address %s:%v", l.Address, err)
			return err
		}
	}
	if m, err := nym.ReceiveResponse(conn); err == nil {
		log.Infof("Received ping message: %s", m)
	}
	return nil
}
