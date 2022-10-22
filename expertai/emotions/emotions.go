// Package emotions provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package emotions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Request
type AnalysisRequest struct {
	// Document
	Document *Document `json:"document,omitempty"`
}

// Categorization data
type CategorizeDocument struct {
	// Categories
	Categories *[]Category `json:"categories,omitempty"`

	// Extra-data containing main groups
	ExtraData *struct {
		// Main groups
		Groups *[]Group `json:"groups,omitempty"`
	} `json:"extraData,omitempty"`

	// Text language
	Language *string `json:"language,omitempty"`

	// Service version
	Version *string `json:"version,omitempty"`
}

// Classification resource response
type CategorizeResponse struct {
	// Categorization data
	Data *CategorizeDocument `json:"data,omitempty"`

	// Service errors
	Errors *[]ServiceError `json:"errors,omitempty"`

	// Operation completed successfully (true/false)
	Success *bool `json:"success,omitempty"`
}

// Category
type Category struct {
	// Score expressed as a percentage of the sum of the scores of all the candidate categories, winners and not (see the **score** property)
	Frequency *float32 `json:"frequency,omitempty"`

	// Hierarchical path
	Hierarchy *[]string `json:"hierarchy,omitempty"`

	// Category ID
	Id *string `json:"id,omitempty"`

	// Category label
	Label *string `json:"label,omitempty"`

	// Name of the software package containing the reference taxonomy
	Namespace *string `json:"namespace,omitempty"`

	// Positions of the portions of text that contributed to the selection of the category
	Positions *[]DocumentPosition `json:"positions,omitempty"`

	// Score assigned to the category to represent its relevance
	Score *int32 `json:"score,omitempty"`

	// True if the category is deemed particularly relevant
	Winner *bool `json:"winner,omitempty"`
}

// Document
type Document struct {
	// The document's text
	Text *string `json:"text,omitempty"`
}

// Position
type DocumentPosition struct {
	// Zero-based position of the character after the last
	End *int64 `json:"end,omitempty"`

	// Zero-based position of the first character
	Start *int64 `json:"start,omitempty"`
}

// Group of emotional traits
type Group struct {
	// ID of the category corresponding to the group inside the taxonomy
	Id *string `json:"id,omitempty"`

	// Label of the category corresponding to the group inside the taxonomy
	Label *string `json:"label,omitempty"`

	// Group rank
	Position *int32 `json:"position,omitempty"`
}

// Error information
type ServiceError struct {
	// Error code
	Code *string `json:"code,omitempty"`

	// Error message
	Message *string `json:"message,omitempty"`
}

// PostCategorizeEmotionalTraitsLanguageJSONBody defines parameters for PostCategorizeEmotionalTraitsLanguage.
type PostCategorizeEmotionalTraitsLanguageJSONBody = AnalysisRequest

// PostCategorizeEmotionalTraitsLanguageParams defines parameters for PostCategorizeEmotionalTraitsLanguage.
type PostCategorizeEmotionalTraitsLanguageParams struct {
	// Classification features, specify `extradata` to obtain main groups
	Features PostCategorizeEmotionalTraitsLanguageParamsFeatures `form:"features" json:"features"`
}

// PostCategorizeEmotionalTraitsLanguageParamsFeatures defines parameters for PostCategorizeEmotionalTraitsLanguage.
type PostCategorizeEmotionalTraitsLanguageParamsFeatures string

// PostCategorizeEmotionalTraitsLanguageParamsLanguage defines parameters for PostCategorizeEmotionalTraitsLanguage.
type PostCategorizeEmotionalTraitsLanguageParamsLanguage string

// PostCategorizeEmotionalTraitsLanguageJSONRequestBody defines body for PostCategorizeEmotionalTraitsLanguage for application/json ContentType.
type PostCategorizeEmotionalTraitsLanguageJSONRequestBody = PostCategorizeEmotionalTraitsLanguageJSONBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// PostCategorizeEmotionalTraitsLanguage request with any body
	PostCategorizeEmotionalTraitsLanguageWithBody(ctx context.Context, language PostCategorizeEmotionalTraitsLanguageParamsLanguage, params *PostCategorizeEmotionalTraitsLanguageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostCategorizeEmotionalTraitsLanguage(ctx context.Context, language PostCategorizeEmotionalTraitsLanguageParamsLanguage, params *PostCategorizeEmotionalTraitsLanguageParams, body PostCategorizeEmotionalTraitsLanguageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) PostCategorizeEmotionalTraitsLanguageWithBody(ctx context.Context, language PostCategorizeEmotionalTraitsLanguageParamsLanguage, params *PostCategorizeEmotionalTraitsLanguageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostCategorizeEmotionalTraitsLanguageRequestWithBody(c.Server, language, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostCategorizeEmotionalTraitsLanguage(ctx context.Context, language PostCategorizeEmotionalTraitsLanguageParamsLanguage, params *PostCategorizeEmotionalTraitsLanguageParams, body PostCategorizeEmotionalTraitsLanguageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostCategorizeEmotionalTraitsLanguageRequest(c.Server, language, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewPostCategorizeEmotionalTraitsLanguageRequest calls the generic PostCategorizeEmotionalTraitsLanguage builder with application/json body
func NewPostCategorizeEmotionalTraitsLanguageRequest(server string, language PostCategorizeEmotionalTraitsLanguageParamsLanguage, params *PostCategorizeEmotionalTraitsLanguageParams, body PostCategorizeEmotionalTraitsLanguageJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostCategorizeEmotionalTraitsLanguageRequestWithBody(server, language, params, "application/json", bodyReader)
}

// NewPostCategorizeEmotionalTraitsLanguageRequestWithBody generates requests for PostCategorizeEmotionalTraitsLanguage with any type of body
func NewPostCategorizeEmotionalTraitsLanguageRequestWithBody(server string, language PostCategorizeEmotionalTraitsLanguageParamsLanguage, params *PostCategorizeEmotionalTraitsLanguageParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "language", runtime.ParamLocationPath, language)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/categorize/emotional-traits/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "features", runtime.ParamLocationQuery, params.Features); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// PostCategorizeEmotionalTraitsLanguage request with any body
	PostCategorizeEmotionalTraitsLanguageWithBodyWithResponse(ctx context.Context, language PostCategorizeEmotionalTraitsLanguageParamsLanguage, params *PostCategorizeEmotionalTraitsLanguageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostCategorizeEmotionalTraitsLanguageResponse, error)

	PostCategorizeEmotionalTraitsLanguageWithResponse(ctx context.Context, language PostCategorizeEmotionalTraitsLanguageParamsLanguage, params *PostCategorizeEmotionalTraitsLanguageParams, body PostCategorizeEmotionalTraitsLanguageJSONRequestBody, reqEditors ...RequestEditorFn) (*PostCategorizeEmotionalTraitsLanguageResponse, error)
}

type PostCategorizeEmotionalTraitsLanguageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *CategorizeResponse
}

// Status returns HTTPResponse.Status
func (r PostCategorizeEmotionalTraitsLanguageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostCategorizeEmotionalTraitsLanguageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// PostCategorizeEmotionalTraitsLanguageWithBodyWithResponse request with arbitrary body returning *PostCategorizeEmotionalTraitsLanguageResponse
func (c *ClientWithResponses) PostCategorizeEmotionalTraitsLanguageWithBodyWithResponse(ctx context.Context, language PostCategorizeEmotionalTraitsLanguageParamsLanguage, params *PostCategorizeEmotionalTraitsLanguageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostCategorizeEmotionalTraitsLanguageResponse, error) {
	rsp, err := c.PostCategorizeEmotionalTraitsLanguageWithBody(ctx, language, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostCategorizeEmotionalTraitsLanguageResponse(rsp)
}

func (c *ClientWithResponses) PostCategorizeEmotionalTraitsLanguageWithResponse(ctx context.Context, language PostCategorizeEmotionalTraitsLanguageParamsLanguage, params *PostCategorizeEmotionalTraitsLanguageParams, body PostCategorizeEmotionalTraitsLanguageJSONRequestBody, reqEditors ...RequestEditorFn) (*PostCategorizeEmotionalTraitsLanguageResponse, error) {
	rsp, err := c.PostCategorizeEmotionalTraitsLanguage(ctx, language, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostCategorizeEmotionalTraitsLanguageResponse(rsp)
}

// ParsePostCategorizeEmotionalTraitsLanguageResponse parses an HTTP response from a PostCategorizeEmotionalTraitsLanguageWithResponse call
func ParsePostCategorizeEmotionalTraitsLanguageResponse(rsp *http.Response) (*PostCategorizeEmotionalTraitsLanguageResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostCategorizeEmotionalTraitsLanguageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest CategorizeResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}
