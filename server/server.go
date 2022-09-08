package server

import (
	"context"
	"net/http"
	"time"

	"github.com/fvbock/endless"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	logging "github.com/ipfs/go-log/v2"

	"github.com/allisterb/citizen5/crypto"
	"github.com/allisterb/citizen5/db"
)

var log = logging.Logger("citizen5/server")

type Config struct {
	Pubkey  string
	PrivKey string
}

func Run(ctx context.Context, config Config) error {
	log.Infof("Starting server...")
	id := crypto.GetIdentity(config.Pubkey)
	log.Infof("server identity is %s.", id.Pretty())
	orbit, dbcleanup, err := db.OpenDB(ctx, config.PrivKey, config.Pubkey)
	if err != nil {
		return err
	}
	db.OpenDocStore(ctx, orbit, "reports")
	if err != nil {
		return err
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(ginzap.Ginzap(log.Desugar(), time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(log.Desugar(), true))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	endless.ListenAndServe(":4242", r)
	orbit.Close()
	dbcleanup()
	return nil
}
