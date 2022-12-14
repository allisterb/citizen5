package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	logging "github.com/ipfs/go-log/v2"

	"github.com/allisterb/citizen5/cmd"
	"github.com/allisterb/citizen5/crypto"
	"github.com/allisterb/citizen5/db"
	"github.com/allisterb/citizen5/nym"
	"github.com/allisterb/citizen5/util"
)

var log = logging.Logger("citizen5/server")

type Config struct {
	Pubkey  string
	PrivKey string
}

func Run(ctx context.Context, config Config, conn *websocket.Conn) error {
	log.Infof("server starting...")
	id := crypto.GetIdentity(config.Pubkey)
	log.Infof("server identity is %s", id.Pretty())
	orbit, dbcleanup, err := db.OpenDB(ctx, config.PrivKey, config.Pubkey)
	if err != nil {
		return err
	}
	reports, err := db.OpenDocStore(ctx, orbit, "reports")
	if err != nil {
		return err
	}

	datastores := db.DataStores{
		DB:      orbit,
		Reports: reports,
	}

	if conn != nil {
		go Monitor(ctx, conn, datastores)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/reports", func(c *gin.Context) {
		allFilter := func(d interface{}) (bool, error) { return true, nil }
		reports, err := datastores.Reports.Query(ctx, allFilter)
		if err != nil {
			c.Status(500)
		} else {
			c.JSON(http.StatusOK, reports)
		}

	})
	r.Use(ginzap.Ginzap(log.Desugar(), time.RFC3339Nano, true))
	r.Use(ginzap.RecoveryWithZap(log.Desugar(), true))
	srv := &http.Server{
		Addr:    ":4242",
		Handler: r,
	}

	go func() {
		log.Infof("REST server started on %s", srv.Addr)
		if err := srv.ListenAndServe(); err == http.ErrServerClosed {
			log.Info("REST server on %s shutdown completed", srv.Addr)
		} else {
			log.Errorf("REST server error: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	log.Infof("server started")
	<-quit
	log.Info("server shutting down...")
	util.Shutdown = true
	srv.Shutdown(ctx)
	orbit.Close()
	dbcleanup()
	log.Info("server shutdown completed")
	return nil
}

func Monitor(ctx context.Context, conn *websocket.Conn, stores db.DataStores) error {
	log.Info("citizen5 Nym service provider running")
	for {
		msg, err := nym.ReceiveMessage(conn)
		if err == nil && util.Shutdown {
			return nil
		} else if err == nil {
			cmd.HandleRemoteCommand(ctx, msg.Binary, stores)
		}
	}
}
