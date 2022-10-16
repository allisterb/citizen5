package nlu

import (
	"context"
	"os"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	logging "github.com/ipfs/go-log/v2"

	"github.com/allisterb/citizen5/expertai/pii"
)

type BearerToken struct {
	LastRefreshed time.Time
	Token         string
}

var log = logging.Logger("citizen5/nlu")
var token = BearerToken{}
var EAIUser = ""
var EAIPass = ""
var PiiClient *pii.Client

func init() {
	EAIUser = os.Getenv("EAI_USER")
	EAIPass = os.Getenv("EAI_PASS")
	c, err := pii.NewClient("https://nlapi.expertapi.")
	if err != nil {
		log.Errorf("Could not create expert.ai PII REST client: %v", err)
		panic("Could not init nlu package.")
	}
	PiiClient = c
}

func RefreshToken() {
	last := time.Since(token.LastRefreshed)

	if last.Hours() > 4 {
		//_, _ := http.Get()
	}
}
func GetPii(ctx context.Context, text string) (string, error) {
	req := pii.PostDetectPiiLanguageJSONRequestBody{}
	req.Document.Text = &text
	bearerAuthProvider, err := securityprovider.NewSecurityProviderBearerToken("MY_USER")
	if err != nil {
		return "", err
	}
	resp, err := PiiClient.PostDetectPiiLanguage(ctx, "en", req, bearerAuthProvider.Intercept)
	return resp.TLS.ServerName, nil
}
