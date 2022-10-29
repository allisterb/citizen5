package cmd

import (
	"context"
	"encoding/json"

	logging "github.com/ipfs/go-log/v2"

	"github.com/allisterb/citizen5/db"
	"github.com/allisterb/citizen5/models"
	"github.com/allisterb/citizen5/nlu"
	"github.com/google/uuid"
)

var log = logging.Logger("citizen5/cmd")

func HandleRemoteCommand(ctx context.Context, cmd []byte, datastores db.DataStores) {
	var doc map[string]interface{}
	var witnessreport models.WitnessReport
	var mediareport models.MediaReport
	if json.Unmarshal(cmd, &witnessreport) == nil {
		log.Infof("received submit witness-report command from %s", witnessreport.Reporter)
		witnessreport.Analysis = models.NLUAnalysis{}
		p, err := nlu.Pii(ctx, witnessreport.Text)
		if err != nil {
			log.Errorf("could not get PII info from expert.ai: %v", err)
		} else {
			witnessreport.Analysis.Pii = p
		}
		if err := json.Unmarshal(cmd, &doc); err != nil {
			log.Errorf("Could not unmarshal WitnessReport data as map")
			return
		}
		doc["_id"] = uuid.New().String()
		_, err = datastores.Reports.Put(ctx, doc)
		if err != nil {
			log.Errorf("error putting witness report to database:%v", err)
		} else {
			log.Infof("report %v stored in database", doc["_id"])
		}
	} else if json.Unmarshal(cmd, &mediareport) == nil {
		log.Infof("received submit media-report command from %s", mediareport.Reporter)
		mediareport.Analysis = models.NLUAnalysis{}
		hs, err := nlu.HateSpeech(ctx, mediareport.Text)
		if err != nil {
			log.Errorf("could not get hate speech info from expert.ai: %v", err)
		} else {
			mediareport.Analysis.HateSpeech = hs
		}
		if err := json.Unmarshal(cmd, &doc); err != nil {
			log.Errorf("could not unmarshal MediaReport data as map")
			return
		}
		doc["_id"] = uuid.New().String()
		_, err = datastores.Reports.Put(ctx, doc)
		if err != nil {
			log.Errorf("error putting media report to database:%v", err)
		} else {
			log.Infof("media report %v stored in database", doc["_id"])
		}
	} else {
		log.Errorf("Unknown command received: %v", string(cmd))
	}
}
