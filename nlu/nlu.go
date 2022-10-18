package nlu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
var Token = BearerToken{}
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

func RefreshToken() error {
	last := time.Since(Token.LastRefreshed)
	if Token.LastRefreshed.IsZero() || last.Hours() > 4 {
		token, err := GetAuthToken()
		if err != nil {
			log.Errorf("Could not refresh expert.ai authorization token: %v", err)
			return err
		}
		Token.Token = token
		Token.LastRefreshed = time.Now()
		return nil
	} else {
		return nil
	}
}

func GetAuthToken() (string, error) {
	data := map[string]string{
		"username": EAIUser,
		"password": EAIPass,
	}
	b, _ := json.Marshal(data)
	bodyReader := bytes.NewReader(b)
	req, _ := http.NewRequest(http.MethodPost, "https://developer.expert.ai/oauth2/token", bodyReader)
	req.Header["Content-Type"] = append(req.Header["Content-Type"], "application/json; charset=utf-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("could not get authorization token from expert.ai API: %v", err)
		return "", err
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("expert.ai API returned status code %v", resp.StatusCode)
		log.Errorf("could not get authorization token from expert.ai API: %v", err)
		return "", nil
	}
	defer resp.Body.Close()
	t, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("could not read authorization token from response: %v", err)
		return "", err
	}
	return string(t), nil
}

func GetPii(ctx context.Context, text string) (string, error) {
	if err := RefreshToken(); err != nil {
		return "", err
	}
	req := pii.PostDetectPiiLanguageJSONRequestBody{}
	req.Document.Text = &text
	bearerAuthProvider, err := securityprovider.NewSecurityProviderBearerToken(Token.Token)
	if err != nil {
		return "", err
	}
	resp, err := PiiClient.PostDetectPiiLanguage(ctx, "en", req, bearerAuthProvider.Intercept)
	return resp.TLS.ServerName, nil
}

func AnalyzeFile(ctx context.Context, file string) error {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Errorf("Could not read file %v", err)
		return err
	}
	GetPii(ctx, string(f))
	return nil
}