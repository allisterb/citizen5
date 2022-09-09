package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/allisterb/citizen5/util"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("citizen5/client")

type Config struct {
	Pubkey  string
	PrivKey string
}

func GetClientConfig() (Config, error) {

	if !util.PathExists(util.ClientConfigFile) {
		log.Errorf("the client config file %s does not exist. Initialize the client first", util.ClientConfigFile)
		return Config{}, fmt.Errorf("client config file not found")
	}
	c, err := ioutil.ReadFile(util.ClientConfigFile)
	if err != nil {
		log.Errorf("Could not read data from client configuration file: %v", err)
		return Config{}, err
	}
	var config Config
	if json.Unmarshal(c, &config) != nil {
		log.Errorf("Could not read JSON data from client configuration file: %v", err)
		return Config{}, err
	} else {
		return config, err
	}
}
