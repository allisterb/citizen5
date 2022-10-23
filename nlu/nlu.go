package nlu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	logging "github.com/ipfs/go-log/v2"

	"github.com/allisterb/citizen5/expertai/nlapi"
	"github.com/allisterb/citizen5/expertai/pii"
)

type BearerToken struct {
	LastRefreshed time.Time
	Token         string
}

var log = logging.Logger("citizen5/nlu")
var Token = BearerToken{}
var EAIUser = os.Getenv("EAI_USER")
var EAIPass = os.Getenv("EAI_PASS")
var PiiClient *pii.Client
var NLApiClient *nlapi.Client

func init() {
	c, err := pii.NewClient("https://nlapi.expert.ai/v2/")
	if err != nil {
		log.Errorf("Could not create expert.ai PII REST client: %v", err)
		panic("Could not init nlu package.")
	}
	PiiClient = c
	a, err := nlapi.NewClient("https://nlapi.expert.ai/v2/")
	if err != nil {
		log.Errorf("Could not create expert.ai NL API REST client: %v", err)
		panic("Could not init nlu package.")
	}
	NLApiClient = a
}

func RefreshToken() error {
	last := time.Since(Token.LastRefreshed)
	if Token.LastRefreshed.IsZero() || last.Hours() > 12 {
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
		return "", err
	}
	defer resp.Body.Close()
	t, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("could not read authorization token from response: %v", err)
		return "", err
	}
	return string(t), nil
}

func Pii(ctx context.Context, text string) (pii.Response, error) {
	var data pii.Response
	if err := RefreshToken(); err != nil {
		return data, err
	}
	req := pii.PostDetectPiiLanguageJSONRequestBody{Document: &struct {
		Text *string "json:\"text,omitempty\""
	}{}}
	req.Document.Text = &text
	bearerAuthProvider, err := securityprovider.NewSecurityProviderBearerToken(Token.Token)
	if err != nil {
		return data, err
	}
	resp, err := PiiClient.PostDetectPiiLanguage(ctx, "en", req, bearerAuthProvider.Intercept)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(b, &data)
	return data, err
}

func Analyze(ctx context.Context, text string) (nlapi.AnalyzeResponse, error) {
	var data nlapi.AnalyzeResponse
	if err := RefreshToken(); err != nil {
		return data, err
	}
	req := nlapi.PostAnalyzeContextLanguageAnalysisJSONBody{Document: &nlapi.Document{Text: &text}}
	req.Document.Text = &text
	bearerAuthProvider, err := securityprovider.NewSecurityProviderBearerToken(Token.Token)
	if err != nil {
		return data, err
	}
	resp, err := NLApiClient.PostAnalyzeContextLanguageAnalysis(ctx, "standard", "en", "", req, bearerAuthProvider.Intercept)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(b, &data)
	return data, err
}
