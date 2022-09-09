package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	orbitdb "berty.tech/go-orbit-db"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	logging "github.com/ipfs/go-log/v2"

	"github.com/allisterb/citizen5/crypto"
	"github.com/allisterb/citizen5/db"
	"github.com/allisterb/citizen5/nym"
)

var log = logging.Logger("citizen5/server")

type Config struct {
	Pubkey  string
	PrivKey string
}

func Run(ctx context.Context, config Config, conn *websocket.Conn) error {
	log.Infof("Starting server...")
	id := crypto.GetIdentity(config.Pubkey)
	log.Infof("server identity is %s.", id.Pretty())
	orbit, dbcleanup, err := db.OpenDB(ctx, config.PrivKey, config.Pubkey)
	if err != nil {
		return err
	}
	reports, err := db.OpenDocStore(ctx, orbit, "reports")
	if err != nil {
		return err
	}
	go Monitor(ctx, conn, reports)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Use(ginzap.Ginzap(log.Desugar(), time.RFC3339Nano, true))
	r.Use(ginzap.RecoveryWithZap(log.Desugar(), true))
	srv := &http.Server{
		Addr:    ":4242",
		Handler: r,
	}

	go func() {
		log.Infof("starting REST server on %s...", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Infof("REST server shutdown requested: %s", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	srv.Shutdown(ctx)
	orbit.Close()
	dbcleanup()
	log.Info("server shutdown completed.")
	return nil
}

func Monitor(ctx context.Context, conn *websocket.Conn, reports orbitdb.DocumentStore) error {
	log.Info("citizen5 Nym service provider running")
	for {
		_, _, err := nym.ReceiveCommand(conn)
		if err != nil {
			return err
		}
	}
}
