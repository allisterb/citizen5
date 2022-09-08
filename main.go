package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	logging "github.com/ipfs/go-log/v2"
	"github.com/mbndr/figlet4go"

	"github.com/allisterb/citizen5/crypto"
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

type ServerCmd struct {
}

var log = logging.Logger("citizen5/main")

// Command-line arguments
var CLI struct {
	Debug      bool          `help:"Enable debug mode."`
	WSUrl      string        `help:"The URL of the Nym websocket client." default:"ws://localhost:1977"`
	Ping       PingCmd       `cmd:"" help:"Send a test message to a Nym mixnet address."`
	Init       InitCmd       `cmd:"" help:"Initialize the citizen5 client."`
	InitServer InitServerCmd `cmd:"" help:"Initialize the citizen5 server."`
	Server     ServerCmd     `cmd:"" help:"Start the citizen5 server."`
}

func init() {
	if os.Getenv("GOLOG_LOG_LEVEL") == "info" { // Reduce noise level of some loggers
		logging.SetLogLevel("dht/RtRefreshManager", "error")
	} else if os.Getenv("GOLOG_LOG_LEVEL") == "" {
		logging.SetAllLoggers(logging.LevelInfo)
		logging.SetLogLevel("dht/RtRefreshManager", "error")
	}
}

func main() {
	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		//figlet4go.ColorBlue,
		figlet4go.ColorCyan,
	}
	renderStr, _ := ascii.RenderOpts("citizenfive", options)
	fmt.Print(renderStr)
	ctx := kong.Parse(&CLI)
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
		if err := nym.SendBinaryFile(conn, nym.GetSelfAddressBinary(conn), "main.go"); err != nil {
			log.Errorf("could not send binary message to Nym mixnet address %s:%v", l.Address, err)
			return err
		}
	} else {
		if err := nym.SendText(conn, l.Address, "ping", true); err != nil {
			log.Errorf("could not send text message to Nym mixnet address %s:%v", l.Address, err)
			return err
		}
	}
	if msg, err := nym.ReceiveMessage(conn); err == nil {
		if len(msg.Binary) == 4 && string(msg.Binary) == "ping" {
			log.Infof("Received ping message OK.")
		} else {
			log.Info(msg.Binary)
		}
	}
	return nil
}

func (l *InitCmd) Run(ctx *kong.Context) error {
	return nil
}

func (s *InitServerCmd) Run(clictx *kong.Context) error {
	ctx, _ := context.WithCancel(context.Background())
	priv, pub := crypto.GenerateIdentity()
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
	log.Infof("IPFS node public key is %s.", crypto.GetIdentity(pub))
	log.Infof("citizen5 server initialized.")
	return nil
}

func (s *ServerCmd) Run(clictx *kong.Context) error {
	if !util.PathExists(util.ServerConfigFile) || !util.PathExists(util.DbDir) {
		log.Errorf("The server config file %s or database directory %s does not exist. Initialize the server first.", util.ServerConfigFile, util.DbDir)
		return fmt.Errorf("server config file or database directory not found")
	}
	c, err := ioutil.ReadFile(util.ServerConfigFile)
	if err != nil {
		log.Errorf("Could not read data from server configuration file: %v", err)
		return err
	}
	var config server.Config
	if json.Unmarshal(c, &config) != nil {
		log.Errorf("Could not read JSON data from server configuration file: %v", err)
		return err
	}
	ctx, _ := context.WithCancel(context.Background())
	conn, err := nym.GetConn(CLI.WSUrl)
	if err != nil {
		log.Errorf("could not open connection to Nym WebSocket %s:%v", CLI.WSUrl, err)
		return nil
	}
	return server.Run(ctx, config, conn)
}
