package server

import (
	"context"

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

	orbit.Close()
	dbcleanup()
	return nil
}
