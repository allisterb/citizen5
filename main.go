package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/alecthomas/kong"
	logging "github.com/ipfs/go-log/v2"
	"github.com/mbndr/figlet4go"

	"github.com/allisterb/citizen5/db"
	"github.com/allisterb/citizen5/nym"
	"github.com/allisterb/citizen5/server"
	"github.com/allisterb/citizen5/util"
)

type PingCmd struct {
	Address string `arg:"" name:"address" help:"Nym mixnet address to send test message to." default:""`
	Binary  bool   `help:"Send a binary file as the test message."`
}

type InitCmd struct {
}

type InitServerCmd struct {
}

var log = logging.Logger("main")

// Command-line arguments
var CLI struct {
	Debug      bool          `help:"Enable debug mode."`
	WSUrl      string        `help:"The URL of the Nym websocket client." default:"ws://localhost:1977"`
	Ping       PingCmd       `cmd:"" help:"Send a test message to a Nym mixnet address."`
	Init       InitCmd       `cmd:"" help:"Initialize the citizen5 client."`
	InitServer InitServerCmd `cmd:"" help:"Initialize the citizen5 server."`
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
	renderStr, _ := ascii.RenderOpts("citizenfive", options)
	fmt.Print(renderStr)
	ctx := kong.Parse(&CLI)
	if util.Contains(ctx.Args, "--debug") {
		logging.SetAllLoggers(logging.LevelDebug)
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

func (s *InitServerCmd) Run(clictx *kong.Context) error {
	ctx, _ := context.WithCancel(context.Background())
	priv, pub := db.GenerateIPFSIdentity()
	err := db.CreateDB(ctx, priv, pub)
	if err != nil {
		log.Errorf("error creating OrbitDB database: %v", err)
		return nil
	}
	serverConfig := server.Config{Pubkey: pub, PrivKey: priv}
	data, _ := json.MarshalIndent(serverConfig, "", " ")
	err = ioutil.WriteFile(filepath.Join(util.GetUserHomeDir(), ".citizen5", "server.json"), data, 0644)
	if err != nil {
		log.Errorf("error creating server configuration file: %v", err)
		return nil
	}
	log.Infof("IPFS node public key is %s", db.GetIPFSIdentity(pub))
	log.Infof("citizen5 server initialized.")
	return nil
}
