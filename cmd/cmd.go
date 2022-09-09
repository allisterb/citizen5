package cmd

import (
	"context"
	"encoding/json"

	logging "github.com/ipfs/go-log/v2"

	"github.com/allisterb/citizen5/db"
	"github.com/allisterb/citizen5/models"
)

var log = logging.Logger("citizen5/cmd")

func HandleRemoteCommand(ctx context.Context, cmd []byte, datastores db.DataStores) {
	var doc map[string]interface{}
	var report models.Report
	if json.Unmarshal(cmd, &report) == nil {
		log.Infof("received submit-report command from %s", report.Id)

		if err := json.Unmarshal(cmd, &doc); err != nil {
			log.Errorf("Could not unmarshal Report data as map")
			return
		}
		datastores.Reports.Put(ctx, doc)
	}
}
