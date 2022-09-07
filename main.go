package main

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
	logging "github.com/ipfs/go-log/v2"
	"github.com/mbndr/figlet4go"

	"github.com/allisterb/citizen5/db"
	"github.com/allisterb/citizen5/nym"
	"github.com/allisterb/citizen5/util"
)

type PingCmd struct {
	Address string `arg:"" name:"address" help:"Nym address to send test message to." default:""`
	Binary  bool   `help:"Send a binary file as the test message."`
}

type InitCmd struct {
}

type CreateDbCmd struct {
	Name string `arg:"" name:"name" help:"Create a citizen5 database with this name."`
}

var log = logging.Logger("main")

// Command-line arguments
var CLI struct {
	Debug bool    `help:"Enable debug mode."`
	WSUrl string  `help:"The URL of the Nym websocket client." default:"ws://localhost:1977"`
	Ping  PingCmd `cmd:"" help:"Send a test message to a Nym mixnet address."`
	Init  InitCmd `cmd:"" help:"Initialze the citizen5 client."`
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
	conn, err := nym.GetConn(CLI.WSUrl)
	if err != nil {
		log.Errorf("could not open connection to Nym WebSocket %s:%v", CLI.WSUrl, err)
		return nil
	}
	defer conn.Close()
	if l.Address == "" {
		l.Address = nym.GetSelfAddressText(conn)
		log.Info("pinging own client address...")
	}
	if l.Binary {
		if err := nym.SendBinary(conn, nym.GetSelfAddressBinary(conn), "main.go"); err != nil {
			log.Errorf("could not send binary message to Nym address %s:%v", l.Address, err)
			return err
		}
	} else {
		if err := nym.SendText(conn, l.Address, "hello", false); err != nil {
			log.Errorf("could not send text message to Nym address %s:%v", l.Address, err)
			return err
		}
	}
	if _, err := nym.ReceiveResponse(conn); err == nil {
		log.Infof("Received ping message OK.")
	}
	return nil
}

func (l *InitCmd) Run(ctx *kong.Context) error {
	return nil
}

func (c *CreateDbCmd) Run(clictx *kong.Context) error {
	ctx, _ := context.WithCancel(context.Background())
	db1, cleanup, err := db.CreateDB(ctx, &c.Name)
	if err != nil {
		log.Errorf("error creating OrbitDB database %s: %v", c.Name, err)
		return nil
	}
	log.Infof("Identity of db: %s", db1.Identity().ID)
	cleanup()
	return nil

}
