package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/allisterb/citizen5/models"
	"github.com/allisterb/citizen5/nlu"
	"github.com/allisterb/citizen5/util"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("citizen5/client")
var Config models.Config

func GetClientConfig() (models.Config, error) {

	if !util.PathExists(util.ClientConfigFile) {
		log.Errorf("the client config file %s does not exist. Initialize the client first", util.ClientConfigFile)
		return models.Config{}, fmt.Errorf("client config file not found")
	}
	c, err := ioutil.ReadFile(util.ClientConfigFile)
	if err != nil {
		log.Errorf("could not read data from client configuration file: %v", err)
		return models.Config{}, err
	}
	var config models.Config
	if json.Unmarshal(c, &config) != nil {
		log.Errorf("could not read JSON data from client configuration file: %v", err)
		return models.Config{}, err
	} else {
		return config, err
	}
}

func SubmitWitnessReport(ctx context.Context, data []byte) error {
	var report models.WitnessReport
	if err := json.Unmarshal(data, &report); err != nil {
		log.Errorf("could not read witness report JSON data from file: %v", err)
		return err
	}
	report.Reporter = Config.Pubkey
	report.DateSubmitted = time.Now().String()
	report.Analysis = models.NLAnalysis{}
	hs, err := nlu.HateSpeech(ctx, report.Description)
	if err != nil {
		log.Errorf("%v", err)
	} else {
		report.Analysis.HateSpeech = hs
	}
	if _, err := json.Marshal(&report); err != nil {
		log.Errorf("Could not create JSON data for submission: %v", err)
		return err
	}
	return nil
}
