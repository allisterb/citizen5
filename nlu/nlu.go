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

	"github.com/allisterb/citizen5/expertai/hatespeech"
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
var HateSpeechClient *hatespeech.Client

func init() {
	c, err := pii.NewClient("https://nlapi.expert.ai/v2/")
	if err != nil {
		log.Errorf("could not create expert.ai PII REST client: %v", err)
		panic("could not init nlu package")
	}
	PiiClient = c
	a, err := nlapi.NewClient("https://nlapi.expert.ai/v2/")
	if err != nil {
		log.Errorf("could not create expert.ai NL API REST client: %v", err)
		panic("could not init nlu package")
	}
	NLApiClient = a
	h, err := hatespeech.NewClient("https://nlapi.expert.ai/v2/")
	if err != nil {
		log.Errorf("could not create expert.ai hate speech REST client: %v", err)
		panic("could not init nlu package")
	}
	HateSpeechClient = h
}

func RefreshToken() error {
	last := time.Since(Token.LastRefreshed)
	if Token.LastRefreshed.IsZero() || last.Hours() > 12 {
		log.Infof("refreshing expert.ai authorization token...")
		token, err := GetAuthToken()
		if err != nil {
			log.Errorf("could not refresh expert.ai authorization token: %v", err)
			return err
		}
		Token.Token = token
		Token.LastRefreshed = time.Now()
		log.Infof("expert.ai authorization token refreshed")
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
	log.Infof("getting expert.ai authorization token...")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("could not get expert.ai authorization token from expert.ai API: %v", err)
		return "", err
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("expert.ai API returned status code %v", resp.StatusCode)
		log.Errorf("could not get expert.ai authorization token from expert.ai API: %v", err)
		return "", err
	}
	defer resp.Body.Close()
	t, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("could not read expert.ai authorization token from response: %v", err)
		return "", err
	}
	log.Infof("get expert.ai authorization token completed")
	return string(t), nil
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

func Pii(ctx context.Context, text string) (pii.Response, error) {
	var data pii.Response
	log.Infof("calling expert.ai PII API...")
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
	log.Infof("calling expert.ai hate speech API completed")
	return data, err
}

func HateSpeech(ctx context.Context, text string) (hatespeech.HateSpeechDetectResponse, error) {
	var data hatespeech.HateSpeechDetectResponse
	log.Infof("calling expert.ai hate speech API...")
	if err := RefreshToken(); err != nil {
		return data, err
	}
	req := hatespeech.PostDetectHateSpeechLanguageJSONRequestBody{Document: &hatespeech.Document{Text: &text}}
	bearerAuthProvider, err := securityprovider.NewSecurityProviderBearerToken(Token.Token)
	if err != nil {
		return data, err
	}
	resp, err := HateSpeechClient.PostDetectHateSpeechLanguage(ctx, "en", req, bearerAuthProvider.Intercept)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(b, &data)
	log.Infof("call expert.ai hate speech API completed")
	return data, err

}

func FindRelation(ctx context.Context, nl nlapi.AnalyzeDocument) error {
	return nil
}
