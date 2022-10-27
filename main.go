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

	"github.com/allisterb/citizen5/client"
	"github.com/allisterb/citizen5/crypto"
	"github.com/allisterb/citizen5/db"
	"github.com/allisterb/citizen5/models"
	"github.com/allisterb/citizen5/nlu"
	"github.com/allisterb/citizen5/nym"
	"github.com/allisterb/citizen5/server"
	"github.com/allisterb/citizen5/util"
)

type PingCmd struct {
	Address string `arg:"" name:"address" help:"Nym mixnet address to send ping message to." default:""`
	Binary  bool   `help:"Send a binary file as the ping message."`
}

type InitCmd struct{}

type InitServerCmd struct{}

type ServerCmd struct{}

type SubmitCmd struct {
	Address string `arg:"" name:"address" help:"The mixnet address of the citizen5 service provider."`
	Type    string `arg:"" name:"type" help:"The type of submission."`
	File    string `arg:"" name:"file" help:"Submit to citizen5 using metadata stored in this file."`
}

type NLUCmd struct {
	File     string `arg:"" name:"file" help:"Analyze text stored in a file."`
	Analysis string `arg:"" name:"analysis" help:"The kind of analysis to perform." default:"full"`
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
	Submit     SubmitCmd     `cmd:"" help:"Submit an item to citizen5."`
	NLU        NLUCmd        `cmd:"" help:"Run NLU models on a plaintext file."`
}

func init() {
	if os.Getenv("GOLOG_LOG_LEVEL") == "info" { // Reduce noise level of some loggers
		logging.SetLogLevel("dht/RtRefreshManager", "error")
		logging.SetLogLevel("bitswap", "error")
		logging.SetLogLevel("connmgr", "error")
	} else if os.Getenv("GOLOG_LOG_LEVEL") == "" {
		logging.SetAllLoggers(logging.LevelInfo)
		logging.SetLogLevel("dht/RtRefreshManager", "error")
		logging.SetLogLevel("bitswap", "error")
		logging.SetLogLevel("connmgr", "error")
		logging.SetLogLevel("net/identify", "error")
	}
}

func main() {
	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
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
		if msg.Json != nil && msg.Json["message"] == "ping" {
			log.Infof("ping %s OK.", l.Address)
		} else {
			log.Errorf("did not receive expected response to ping from %s.", l.Address)
		}
	}
	return nil
}

func (c *InitCmd) Run(clictx *kong.Context) error {
	priv, pub := crypto.GenerateIdentity()
	clientConfig := models.Config{Pubkey: pub, PrivKey: priv}
	data, _ := json.MarshalIndent(clientConfig, "", " ")
	err := ioutil.WriteFile(filepath.Join(util.GetUserHomeDir(), ".citizen5", "client.json"), data, 0644)
	if err != nil {
		log.Errorf("error creating client configuration file: %v", err)
		return nil
	}
	log.Infof("client identity is %s.", crypto.GetIdentity(pub).Pretty())
	log.Infof("citizen5 client initialized.")
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
	log.Infof("server identity is %s.", crypto.GetIdentity(pub).Pretty())
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
		log.Errorf("could not read data from server configuration file: %v", err)
		return err
	}
	var config server.Config
	if json.Unmarshal(c, &config) != nil {
		log.Errorf("could not read JSON data from server configuration file: %v", err)
		return err
	}
	ctx, _ := context.WithCancel(context.Background())
	conn, err := nym.GetConn(CLI.WSUrl)
	if err != nil {
		log.Errorf("could not open connection to Nym WebSocket %s:%v", CLI.WSUrl, err)
		return nil
	}
	defer conn.Close()
	return server.Run(ctx, config, conn)
}

func (r *SubmitCmd) Run(clictx *kong.Context) error {
	switch r.Type {
	case "report", "witness-report", "media-report":
		break
	default:
		err := fmt.Errorf("unknown submission type: %v", r.Type)
		return err
	}
	config, err := client.GetClientConfig()
	if err != nil {
		return nil
	}
	client.Config = config
	if !util.PathExists(r.File) {
		err = fmt.Errorf("the file %s does not exist", r.File)
		return err
	}
	c, err := ioutil.ReadFile(r.File)
	if err != nil {
		log.Errorf("could not read data from file: %v", err)
		return err
	}
	ctx, _ := context.WithCancel(context.Background())
	conn, err := nym.GetConn(CLI.WSUrl)
	if err != nil {
		log.Errorf("could not open connection to Nym WebSocket %s:%v", CLI.WSUrl, err)
		return err
	}
	switch r.Type {
	case "witness-report":
		return client.SubmitWitnessReport(ctx, conn, r.Address, c)
	case "media-report":
		return client.SubmitMediaReport(ctx, conn, r.Address, c)
	default:
		panic(fmt.Errorf("unknown submission type: %v", r.Type))
	}
}

func (c *NLUCmd) Run(clictx *kong.Context) error {
	ctx, _ := context.WithCancel(context.Background())
	f, err := ioutil.ReadFile(c.File)
	if err != nil {
		log.Errorf("could not read file %v", err)
		return err
	}
	switch c.Analysis {
	case "pii":
		p, err := nlu.Pii(ctx, string(f))
		if p.Success == nil || !*p.Success {
			log.Errorf("could not get PII from expert.ai API for file %v: %v", c.File, err)
			return err
		}
		j, _ := json.MarshalIndent(p.Data.Extractions, "", "  ")
		log.Info(string(j))

	case "topics", "relations", "entities", "lemmas", "mainphrases", "full":
		p, err := nlu.Analyze(ctx, string(f))
		if p.Success == nil || !*p.Success {
			log.Errorf("could not get call expert.ai NLU API for file %v: %v", c.File, err)
			return err
		}
		switch c.Analysis {
		case "topics":
			j, _ := json.MarshalIndent(p.Data.Topics, "", "  ")
			log.Infof("printing topics in %v", c.File)
			log.Info(string(j))
		case "relations":
			j, _ := json.MarshalIndent(p.Data.Relations, "", "  ")
			log.Infof("printing relations in %v", c.File)
			log.Info(string(j))
		case "entities":
			j, _ := json.MarshalIndent(p.Data.Entities, "", "  ")
			log.Infof("printing entities in %v", c.File)
			log.Info(string(j))
		case "lemmas":
			j, _ := json.MarshalIndent(p.Data.MainLemmas, "", "  ")
			log.Infof("printing main lemmas in %v", c.File)
			log.Info(string(j))
		case "mainphrases":
			j, _ := json.MarshalIndent(p.Data.MainPhrases, "", "  ")
			log.Infof("printing main phrases in %v", c.File)
			log.Info(string(j))
		case "full":
			j, _ := json.MarshalIndent(p.Data, "", "  ")
			log.Infof("printing full analysis for %v...", c.File)
			log.Info(string(j))
		}
	case "hatespeech":
		p, err := nlu.HateSpeech(ctx, string(f))
		if p.Success == nil || !*p.Success {
			log.Errorf("Could not hate speech analysis from expert.ai API for file %v: %v", c.File, err)
			return err
		}
		j, _ := json.MarshalIndent(p.Data, "", "  ")
		log.Infof("printing hate speech analysis for %v...", c.File)
		log.Info(string(j))

	default:
		err := fmt.Errorf("unknown analysis: %v", c.Analysis)
		log.Errorf("%v", err)
		return err
	}
	return nil
}
