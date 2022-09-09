package cmd

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/allisterb/citizen5/db"
	"github.com/allisterb/citizen5/models"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("citizen5/cmd")

func HandleRemoteCommand(ctx context.Context, cmd []byte, datastores db.DataStores) {
	var doc map[string]interface{}
	var report models.Report
	if json.Unmarshal(cmd, &report) == nil {
		log.Infof("received submit-report command from %s", report.Reporter)

		if err := json.Unmarshal(cmd, &doc); err != nil {
			log.Errorf("Could not unmarshal Report data as map")
			return
		}
		doc["_id"] = uuid.New().String()
		_, err := datastores.Reports.Put(ctx, doc)
		if err != nil {
			log.Errorf("error putting report to database:%v", err)
		}
		//err := datastores.Reports.Sync()
		datastores.DB.Close()
	}
}
