// Package pii provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package pii

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

// Atom info
type Atom struct {
	// Zero-based position of the first character after the atom
	End *int64 `json:"end,omitempty"`

	// Lemma
	Lemma *string `json:"lemma,omitempty"`

	// Zero-based position of the first character of the atom
	Start *int64 `json:"start,omitempty"`

	// expert.ai type
	Type *string `json:"type,omitempty"`
}

// Dependency info
type Dependency struct {
	// Number of the head token
	Head *int64 `json:"head,omitempty"`

	// Zero-based cardinal number of the token
	Id *int64 `json:"id,omitempty"`

	// <a href='https://universaldependencies.org/u/dep/#universal-dependency-relations' target='_blank'>Dependency relation</a> between the token and the head token
	Label *string `json:"label,omitempty"`
}

// Position
type DocumentPosition struct {
	// Zero-based position of the character after the last
	End *int64 `json:"end,omitempty"`

	// Zero-based position of the first character
	Start *int64 `json:"start,omitempty"`
}

// Entity info
type Entity struct {
	// Entity attributes inferred from the context or from the Knowledge Graph
	Attributes *[]InferredAttribute `json:"attributes,omitempty"`

	// Base form (lemma) of the entity name
	Lemma *string `json:"lemma,omitempty"`

	// Positions of the entity's mentions
	Positions *[]DocumentPosition `json:"positions,omitempty"`

	// Entity relevance
	Relevance *int64 `json:"relevance,omitempty"`

	// ID used to look up Knowledge Graph data in the `knowledge` array
	Syncon *int64 `json:"syncon,omitempty"`

	// Entity type
	Type *string `json:"type,omitempty"`
}

// Extraction record
type Extraction struct {
	// Extraction record fields
	Fields *[]struct {
		// Field name
		Name *string `json:"name,omitempty"`

		// Positions of parts of the text corresponding to the field value
		Positions *[]DocumentPosition `json:"positions,omitempty"`

		// Field value
		Value *string `json:"value,omitempty"`
	} `json:"fields,omitempty"`

	// Software package name
	Namespace *string `json:"namespace,omitempty"`

	// Extraction record template
	Template *string `json:"template,omitempty"`
}

// Inferred attribute
type InferredAttribute struct {
	// Attribute name
	Attribute *string `json:"attribute,omitempty"`

	// Attribute's attributes
	Attributes *[]InferredAttribute `json:"attributes,omitempty"`

	// Lemma
	Lemma *string `json:"lemma,omitempty"`

	// ID used to look up Knowledge Graph data in the `knowledge` array
	Syncon *int64 `json:"syncon,omitempty"`

	// Attribute type
	Type *string `json:"type,omitempty"`
}

// Knowledge Graph data for a syncon
type KnowledgeEntry struct {
	// Textual rendering of the general conceptual category for the token in the Knowledge Graph
	Label *string `json:"label,omitempty"`

	// Syncon extended properties
	Properties *[]Property `json:"properties,omitempty"`

	// Syncon ID
	Syncon *int64 `json:"syncon,omitempty"`
}

// PIIBankAccount defines model for PIIBankAccount.
type PIIBankAccount struct {
	// JSON-LD graph item id
	Id *string `json:"@id,omitempty"`

	// JSON-LD graph item type
	Type        *string `json:"@type,omitempty"`
	IBAN        *string `json:"IBAN,omitempty"`
	IBANcountry *string `json:"IBANcountry,omitempty"`

	// Text matches of items' properties
	Matches *[]struct {
		// Property value in the text, zero-based position of the character after the last
		End *int64 `json:"end,omitempty"`

		// Property name
		Name *string `json:"name,omitempty"`

		// Property value in the text, zero-based position of the first character
		Start *int64 `json:"start,omitempty"`

		// Property value
		Value *string `json:"value,omitempty"`
	} `json:"matches,omitempty"`
}

// Type with properties inherited by all JSON-LD graph PII types
type PIIBaseItem struct {
	// JSON-LD graph item id
	Id *string `json:"@id,omitempty"`

	// JSON-LD graph item type
	Type *string `json:"@type,omitempty"`

	// Text matches of items' properties
	Matches *[]struct {
		// Property value in the text, zero-based position of the character after the last
		End *int64 `json:"end,omitempty"`

		// Property name
		Name *string `json:"name,omitempty"`

		// Property value in the text, zero-based position of the first character
		Start *int64 `json:"start,omitempty"`

		// Property value
		Value *string `json:"value,omitempty"`
	} `json:"matches,omitempty"`
}

// PIIFinancialProduct defines model for PIIFinancialProduct.
type PIIFinancialProduct struct {
	// JSON-LD graph item id
	Id *string `json:"@id,omitempty"`

	// JSON-LD graph item type
	Type              *string `json:"@type,omitempty"`
	CVV               *string `json:"CVV,omitempty"`
	CreditDebitNumber *string `json:"creditDebitNumber,omitempty"`
	ExpirationDate    *string `json:"expirationDate,omitempty"`

	// Text matches of items' properties
	Matches *[]struct {
		// Property value in the text, zero-based position of the character after the last
		End *int64 `json:"end,omitempty"`

		// Property name
		Name *string `json:"name,omitempty"`

		// Property value in the text, zero-based position of the first character
		Start *int64 `json:"start,omitempty"`

		// Property value
		Value *string `json:"value,omitempty"`
	} `json:"matches,omitempty"`
}

// PIIIP defines model for PIIIP.
type PIIIP struct {
	// JSON-LD graph item id
	Id *string `json:"@id,omitempty"`

	// JSON-LD graph item type
	Type *string `json:"@type,omitempty"`
	IP   *string `json:"IP,omitempty"`

	// Text matches of items' properties
	Matches *[]struct {
		// Property value in the text, zero-based position of the character after the last
		End *int64 `json:"end,omitempty"`

		// Property name
		Name *string `json:"name,omitempty"`

		// Property value in the text, zero-based position of the first character
		Start *int64 `json:"start,omitempty"`

		// Property value
		Value *string `json:"value,omitempty"`
	} `json:"matches,omitempty"`
}

// Personally Identifiable Information (PII) item
type PIIItem interface{}

// PIIPerson defines model for PIIPerson.
type PIIPerson struct {
	// JSON-LD graph item id
	Id *string `json:"@id,omitempty"`

	// JSON-LD graph item type
	Type       *string   `json:"@type,omitempty"`
	Age        *string   `json:"age,omitempty"`
	BirthDate  *string   `json:"birthDate,omitempty"`
	BirthPlace *string   `json:"birthPlace,omitempty"`
	DateTime   *[]string `json:"dateTime,omitempty"`
	DeathDate  *string   `json:"deathDate,omitempty"`
	DeathPlace *string   `json:"deathPlace,omitempty"`
	FamilyName *string   `json:"familyName,omitempty"`
	Gender     *string   `json:"gender,omitempty"`
	GivenName  *string   `json:"givenName,omitempty"`

	// Text matches of items' properties
	Matches *[]struct {
		// Property value in the text, zero-based position of the character after the last
		End *int64 `json:"end,omitempty"`

		// Property name
		Name *string `json:"name,omitempty"`

		// Property value in the text, zero-based position of the first character
		Start *int64 `json:"start,omitempty"`

		// Property value
		Value *string `json:"value,omitempty"`
	} `json:"matches,omitempty"`
	Nationality *string `json:"nationality,omitempty"`
	Person      *string `json:"person,omitempty"`
}

// PIIPostalAddress defines model for PIIPostalAddress.
type PIIPostalAddress struct {
	// JSON-LD graph item id
	Id *string `json:"@id,omitempty"`

	// JSON-LD graph item type
	Type            *string `json:"@type,omitempty"`
	Address         *string `json:"address,omitempty"`
	AddressCountry  *string `json:"addressCountry,omitempty"`
	AddressLocality *string `json:"addressLocality,omitempty"`
	AddressRegion   *string `json:"addressRegion,omitempty"`

	// Text matches of items' properties
	Matches *[]struct {
		// Property value in the text, zero-based position of the character after the last
		End *int64 `json:"end,omitempty"`

		// Property name
		Name *string `json:"name,omitempty"`

		// Property value in the text, zero-based position of the first character
		Start *int64 `json:"start,omitempty"`

		// Property value
		Value *string `json:"value,omitempty"`
	} `json:"matches,omitempty"`
	PostOfficeBoxNumber *string `json:"postOfficeBoxNumber,omitempty"`
	PostalCode          *string `json:"postalCode,omitempty"`
	StreetAddress       *string `json:"streetAddress,omitempty"`
}

// PIIURL defines model for PIIURL.
type PIIURL struct {
	// JSON-LD graph item id
	Id *string `json:"@id,omitempty"`

	// JSON-LD graph item type
	Type *string `json:"@type,omitempty"`
	URL  *string `json:"URL,omitempty"`

	// Text matches of items' properties
	Matches *[]struct {
		// Property value in the text, zero-based position of the character after the last
		End *int64 `json:"end,omitempty"`

		// Property name
		Name *string `json:"name,omitempty"`

		// Property value in the text, zero-based position of the first character
		Start *int64 `json:"start,omitempty"`

		// Property value
		Value *string `json:"value,omitempty"`
	} `json:"matches,omitempty"`
}

// PIIemail defines model for PIIemail.
type PIIemail struct {
	// JSON-LD graph item id
	Id *string `json:"@id,omitempty"`

	// JSON-LD graph item type
	Type  *string `json:"@type,omitempty"`
	Email *string `json:"email,omitempty"`

	// Text matches of items' properties
	Matches *[]struct {
		// Property value in the text, zero-based position of the character after the last
		End *int64 `json:"end,omitempty"`

		// Property name
		Name *string `json:"name,omitempty"`

		// Property value in the text, zero-based position of the first character
		Start *int64 `json:"start,omitempty"`

		// Property value
		Value *string `json:"value,omitempty"`
	} `json:"matches,omitempty"`
}

// PIItelephone defines model for PIItelephone.
type PIItelephone struct {
	// JSON-LD graph item id
	Id *string `json:"@id,omitempty"`

	// JSON-LD graph item type
	Type *string `json:"@type,omitempty"`

	// Text matches of items' properties
	Matches *[]struct {
		// Property value in the text, zero-based position of the character after the last
		End *int64 `json:"end,omitempty"`

		// Property name
		Name *string `json:"name,omitempty"`

		// Property value in the text, zero-based position of the first character
		Start *int64 `json:"start,omitempty"`

		// Property value
		Value *string `json:"value,omitempty"`
	} `json:"matches,omitempty"`
	Telephone *string `json:"telephone,omitempty"`
}

// Paragraph info
type Paragraph struct {
	// Zero-based position of the first character after the paragraph
	End *int64 `json:"end,omitempty"`

	// Indexes (in the `sentences` array) of the sentences that make up the paragraph
	Sentences *[]int64 `json:"sentences,omitempty"`

	// Zero-based position of the first character of the paragraph
	Start *int64 `json:"start,omitempty"`
}

// Phrase info
type Phrase struct {
	// Zero-based position of the first character after the phrase
	End *int64 `json:"end,omitempty"`

	// Zero-based position of the first character of the phrase
	Start *int64 `json:"start,omitempty"`

	// Indexes (in the `tokens` array) of the tokens that make up the phrase
	Tokens *[]int64 `json:"tokens,omitempty"`

	// Phrase type
	Type *string `json:"type,omitempty"`
}

// Syncon property
type Property struct {
	// Property type
	Type *string `json:"type,omitempty"`

	// Property value
	Value *string `json:"value,omitempty"`
}

// Request
type Request struct {
	// Document to analyze
	Document *struct {
		// Document's text
		Text *string `json:"text,omitempty"`
	} `json:"document,omitempty"`
}

// Detector's response
type Response struct {
	Data *struct {
		// Analyzed text
		Content *string `json:"content,omitempty"`

		// Entities
		Entities *[]Entity `json:"entities,omitempty"`

		// PII detector specific output
		ExtraData *struct {
			// JSON-LD format output
			JSONLD *struct {
				// JSON-LD context
				Context *struct {
					Version            *string `json:"@version,omitempty"`
					CVV                *string `json:"CVV,omitempty"`
					IBAN               *string `json:"IBAN,omitempty"`
					IBANcountry        *string `json:"IBANcountry,omitempty"`
					IP                 *string `json:"IP,omitempty"`
					URL                *string `json:"URL,omitempty"`
					AdditionalProperty *string `json:"additionalProperty,omitempty"`
					Address            *string `json:"address,omitempty"`
					AddressCountry     *string `json:"addressCountry,omitempty"`
					AddressLocality    *string `json:"addressLocality,omitempty"`
					AddressRegion      *string `json:"addressRegion,omitempty"`
					Age                *string `json:"age,omitempty"`
					BirthDate          *string `json:"birthDate,omitempty"`
					BirthPlace         *string `json:"birthPlace,omitempty"`
					CreditDebitNumber  *string `json:"creditDebitNumber,omitempty"`
					DateTime           *string `json:"dateTime,omitempty"`
					DeathDate          *string `json:"deathDate,omitempty"`
					DeathPlace         *string `json:"deathPlace,omitempty"`
					Email              *string `json:"email,omitempty"`
					End                *struct {
						Id *string `json:"@id,omitempty"`
					} `json:"end,omitempty"`
					ExpirationDate      *string `json:"expirationDate,omitempty"`
					FamilyName          *string `json:"familyName,omitempty"`
					Gender              *string `json:"gender,omitempty"`
					GivenName           *string `json:"givenName,omitempty"`
					Matches             *string `json:"matches,omitempty"`
					Nationality         *string `json:"nationality,omitempty"`
					Person              *string `json:"person,omitempty"`
					PostOfficeBoxNumber *string `json:"postOfficeBoxNumber,omitempty"`
					PostalAddress       *string `json:"postalAddress,omitempty"`
					PostalCode          *string `json:"postalCode,omitempty"`
					Schema              *string `json:"schema,omitempty"`
					Start               *struct {
						Id *string `json:"@id,omitempty"`
					} `json:"start,omitempty"`
					StreetAddress *string `json:"streetAddress,omitempty"`
					Telephone     *string `json:"telephone,omitempty"`
					Text          *string `json:"text,omitempty"`
					Type          *string `json:"type,omitempty"`
				} `json:"@context,omitempty"`

				// JSON-LD graph data
				Graph *[]PIIItem `json:"@graph,omitempty"`
			} `json:"JSON-LD,omitempty"`
		} `json:"extraData,omitempty"`

		// Extracted information
		Extractions *[]Extraction `json:"extractions,omitempty"`

		// Knowledge Graph syncons' data
		Knowledge *[]KnowledgeEntry `json:"knowledge,omitempty"`

		// Text language
		Language *string `json:"language,omitempty"`

		// Paragraphs
		Paragraphs *[]Paragraph `json:"paragraphs,omitempty"`

		// Phrases
		Phrases *[]Phrase `json:"phrases,omitempty"`

		// Sentences
		Sentences *[]Sentence `json:"sentences,omitempty"`

		// Tokens
		Tokens *[]Token `json:"tokens,omitempty"`

		// Service version
		Version *string `json:"version,omitempty"`
	} `json:"data,omitempty"`

	// Error description
	Error *struct {
		// Error code
		Code *string `json:"code,omitempty"`

		// Error message
		Message *string `json:"message,omitempty"`
	} `json:"error,omitempty"`

	// Success flag
	Success *bool `json:"success,omitempty"`
}

// Sentence info
type Sentence struct {
	// Zero-based position of the first character after the sentence
	End *int64 `json:"end,omitempty"`

	// Indexes (in the `phrases` array) of the phrases that make up the sentence
	Phrases *[]int64 `json:"phrases,omitempty"`

	// Zero-based position of the first character of the sentence
	Start *int64 `json:"start,omitempty"`
}

// Token info
type Token struct {
	// Atoms that make up the token
	Atoms *[]Atom `json:"atoms,omitempty"`

	// Dependency info
	Dependency *Dependency `json:"dependency,omitempty"`

	// Zero-based position of the first character after the token
	End *int64 `json:"end,omitempty"`

	// Lemma
	Lemma *string `json:"lemma,omitempty"`

	// A semicolon separated list of <a href='https://universaldependencies.org/format.html#morphological-annotation'>CoNLL-U format</a> morphological features
	Morphology *string `json:"morphology,omitempty"`

	// Paragraph index in the `paragraphs` array
	Paragraph *int64 `json:"paragraph,omitempty"`

	// Phrase index in the `phrases` array
	Phrase *int64 `json:"phrase,omitempty"`

	// <a href='https://universaldependencies.org/u/pos/'>Universal Dependencies part-of-speech tag</a>
	Pos *string `json:"pos,omitempty"`

	// Sentence index in the `sentences` array
	Sentence *int64 `json:"sentence,omitempty"`

	// Zero-based position of the first character of the token
	Start *int64 `json:"start,omitempty"`

	// ID used to look up Knowledge Graph data in the `knowledge` array
	Syncon *int64 `json:"syncon,omitempty"`

	// expert.ai type
	Type *string `json:"type,omitempty"`

	// A concept that does not exist in the Knowledge Graph but heuristics recognized as a type of a known parent concept
	Vsyn *VirtualSyncon `json:"vsyn,omitempty"`
}

// A concept that does not exist in the Knowledge Graph but heuristics recognized as a type of a known parent concept
type VirtualSyncon struct {
	// ID used to mark all the occurrences of the virtual concept in the text
	Id *int64 `json:"id,omitempty"`

	// Parent concept; ID is used to look up Knowledge Graph data in the `knowledge` array
	Parent *int64 `json:"parent,omitempty"`
}

// PostDetectPiiLanguageJSONBody defines parameters for PostDetectPiiLanguage.
type PostDetectPiiLanguageJSONBody = Request

// PostDetectPiiLanguageJSONRequestBody defines body for PostDetectPiiLanguage for application/json ContentType.
type PostDetectPiiLanguageJSONRequestBody = PostDetectPiiLanguageJSONBody

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
	// PostDetectPiiLanguage request with any body
	PostDetectPiiLanguageWithBody(ctx context.Context, language string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostDetectPiiLanguage(ctx context.Context, language string, body PostDetectPiiLanguageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) PostDetectPiiLanguageWithBody(ctx context.Context, language string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostDetectPiiLanguageRequestWithBody(c.Server, language, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostDetectPiiLanguage(ctx context.Context, language string, body PostDetectPiiLanguageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostDetectPiiLanguageRequest(c.Server, language, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewPostDetectPiiLanguageRequest calls the generic PostDetectPiiLanguage builder with application/json body
func NewPostDetectPiiLanguageRequest(server string, language string, body PostDetectPiiLanguageJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostDetectPiiLanguageRequestWithBody(server, language, "application/json", bodyReader)
}

// NewPostDetectPiiLanguageRequestWithBody generates requests for PostDetectPiiLanguage with any type of body
func NewPostDetectPiiLanguageRequestWithBody(server string, language string, contentType string, body io.Reader) (*http.Request, error) {
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

	operationPath := fmt.Sprintf("/detect/pii/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

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
	// PostDetectPiiLanguage request with any body
	PostDetectPiiLanguageWithBodyWithResponse(ctx context.Context, language string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostDetectPiiLanguageResponse, error)

	PostDetectPiiLanguageWithResponse(ctx context.Context, language string, body PostDetectPiiLanguageJSONRequestBody, reqEditors ...RequestEditorFn) (*PostDetectPiiLanguageResponse, error)
}

type PostDetectPiiLanguageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Response
}

// Status returns HTTPResponse.Status
func (r PostDetectPiiLanguageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostDetectPiiLanguageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// PostDetectPiiLanguageWithBodyWithResponse request with arbitrary body returning *PostDetectPiiLanguageResponse
func (c *ClientWithResponses) PostDetectPiiLanguageWithBodyWithResponse(ctx context.Context, language string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostDetectPiiLanguageResponse, error) {
	rsp, err := c.PostDetectPiiLanguageWithBody(ctx, language, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostDetectPiiLanguageResponse(rsp)
}

func (c *ClientWithResponses) PostDetectPiiLanguageWithResponse(ctx context.Context, language string, body PostDetectPiiLanguageJSONRequestBody, reqEditors ...RequestEditorFn) (*PostDetectPiiLanguageResponse, error) {
	rsp, err := c.PostDetectPiiLanguage(ctx, language, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostDetectPiiLanguageResponse(rsp)
}

// ParsePostDetectPiiLanguageResponse parses an HTTP response from a PostDetectPiiLanguageWithResponse call
func ParsePostDetectPiiLanguageResponse(rsp *http.Response) (*PostDetectPiiLanguageResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostDetectPiiLanguageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Response
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}
