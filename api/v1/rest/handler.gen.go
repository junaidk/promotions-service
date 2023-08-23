// Package rest provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
package rest

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// Error defines model for Error.
type Error struct {
	Error  *string `json:"error,omitempty"`
	Status *string `json:"status,omitempty"`
}

// ProcessCsv defines model for ProcessCsv.
type ProcessCsv struct {
	FilePath *string `json:"file_path,omitempty"`
}

// Promotion defines model for Promotion.
type Promotion struct {
	ExpirationDate *string `json:"expiration_date,omitempty"`
	Id             *string `json:"id,omitempty"`
	Price          *string `json:"price,omitempty"`
}

// PostV1AdminProcessCsvJSONRequestBody defines body for PostV1AdminProcessCsv for application/json ContentType.
type PostV1AdminProcessCsvJSONRequestBody = ProcessCsv

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
	// PostV1AdminProcessCsv request with any body
	PostV1AdminProcessCsvWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostV1AdminProcessCsv(ctx context.Context, body PostV1AdminProcessCsvJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PostV1AdminSwitchDb request
	PostV1AdminSwitchDb(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetV1PromotionsId request
	GetV1PromotionsId(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) PostV1AdminProcessCsvWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostV1AdminProcessCsvRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostV1AdminProcessCsv(ctx context.Context, body PostV1AdminProcessCsvJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostV1AdminProcessCsvRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostV1AdminSwitchDb(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostV1AdminSwitchDbRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetV1PromotionsId(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetV1PromotionsIdRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewPostV1AdminProcessCsvRequest calls the generic PostV1AdminProcessCsv builder with application/json body
func NewPostV1AdminProcessCsvRequest(server string, body PostV1AdminProcessCsvJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostV1AdminProcessCsvRequestWithBody(server, "application/json", bodyReader)
}

// NewPostV1AdminProcessCsvRequestWithBody generates requests for PostV1AdminProcessCsv with any type of body
func NewPostV1AdminProcessCsvRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/admin/process-csv")
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

// NewPostV1AdminSwitchDbRequest generates requests for PostV1AdminSwitchDb
func NewPostV1AdminSwitchDbRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/admin/switch-db")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetV1PromotionsIdRequest generates requests for GetV1PromotionsId
func NewGetV1PromotionsIdRequest(server string, id string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "id", runtime.ParamLocationPath, id)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/promotions/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

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
	// PostV1AdminProcessCsv request with any body
	PostV1AdminProcessCsvWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostV1AdminProcessCsvResponse, error)

	PostV1AdminProcessCsvWithResponse(ctx context.Context, body PostV1AdminProcessCsvJSONRequestBody, reqEditors ...RequestEditorFn) (*PostV1AdminProcessCsvResponse, error)

	// PostV1AdminSwitchDb request
	PostV1AdminSwitchDbWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PostV1AdminSwitchDbResponse, error)

	// GetV1PromotionsId request
	GetV1PromotionsIdWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*GetV1PromotionsIdResponse, error)
}

type PostV1AdminProcessCsvResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r PostV1AdminProcessCsvResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostV1AdminProcessCsvResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostV1AdminSwitchDbResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r PostV1AdminSwitchDbResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostV1AdminSwitchDbResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetV1PromotionsIdResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Promotion
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetV1PromotionsIdResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetV1PromotionsIdResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// PostV1AdminProcessCsvWithBodyWithResponse request with arbitrary body returning *PostV1AdminProcessCsvResponse
func (c *ClientWithResponses) PostV1AdminProcessCsvWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostV1AdminProcessCsvResponse, error) {
	rsp, err := c.PostV1AdminProcessCsvWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostV1AdminProcessCsvResponse(rsp)
}

func (c *ClientWithResponses) PostV1AdminProcessCsvWithResponse(ctx context.Context, body PostV1AdminProcessCsvJSONRequestBody, reqEditors ...RequestEditorFn) (*PostV1AdminProcessCsvResponse, error) {
	rsp, err := c.PostV1AdminProcessCsv(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostV1AdminProcessCsvResponse(rsp)
}

// PostV1AdminSwitchDbWithResponse request returning *PostV1AdminSwitchDbResponse
func (c *ClientWithResponses) PostV1AdminSwitchDbWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PostV1AdminSwitchDbResponse, error) {
	rsp, err := c.PostV1AdminSwitchDb(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostV1AdminSwitchDbResponse(rsp)
}

// GetV1PromotionsIdWithResponse request returning *GetV1PromotionsIdResponse
func (c *ClientWithResponses) GetV1PromotionsIdWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*GetV1PromotionsIdResponse, error) {
	rsp, err := c.GetV1PromotionsId(ctx, id, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetV1PromotionsIdResponse(rsp)
}

// ParsePostV1AdminProcessCsvResponse parses an HTTP response from a PostV1AdminProcessCsvWithResponse call
func ParsePostV1AdminProcessCsvResponse(rsp *http.Response) (*PostV1AdminProcessCsvResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostV1AdminProcessCsvResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParsePostV1AdminSwitchDbResponse parses an HTTP response from a PostV1AdminSwitchDbWithResponse call
func ParsePostV1AdminSwitchDbResponse(rsp *http.Response) (*PostV1AdminSwitchDbResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostV1AdminSwitchDbResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGetV1PromotionsIdResponse parses an HTTP response from a GetV1PromotionsIdWithResponse call
func ParseGetV1PromotionsIdResponse(rsp *http.Response) (*GetV1PromotionsIdResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetV1PromotionsIdResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Promotion
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /v1/admin/process-csv)
	PostV1AdminProcessCsv(w http.ResponseWriter, r *http.Request)

	// (POST /v1/admin/switch-db)
	PostV1AdminSwitchDb(w http.ResponseWriter, r *http.Request)

	// (GET /v1/promotions/{id})
	GetV1PromotionsId(w http.ResponseWriter, r *http.Request, id string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// PostV1AdminProcessCsv operation middleware
func (siw *ServerInterfaceWrapper) PostV1AdminProcessCsv(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostV1AdminProcessCsv(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostV1AdminSwitchDb operation middleware
func (siw *ServerInterfaceWrapper) PostV1AdminSwitchDb(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostV1AdminSwitchDb(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetV1PromotionsId operation middleware
func (siw *ServerInterfaceWrapper) GetV1PromotionsId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, chi.URLParam(r, "id"), &id)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetV1PromotionsId(w, r, id)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshallingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshallingParamError) Error() string {
	return fmt.Sprintf("Error unmarshalling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshallingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/admin/process-csv", wrapper.PostV1AdminProcessCsv)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/admin/switch-db", wrapper.PostV1AdminSwitchDb)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1/promotions/{id}", wrapper.GetV1PromotionsId)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9xUTYvbMBD9K2bao3ftdLcX39IPSi5lYWEvbVgUaZIoxJI6Gqddgv97GSlOwsaHHBYK",
	"PVmM5uPpvefZg/Zt8A4dR2j2EPUaW5WOX4k8ySGQD0hsMYVxCPNLQGggMlm3gr6EyIq7OHLVl0PELzao",
	"WZIfyGuM8XPcXY5Y2i0+B8Xr63u1nq13I2j/BEtK7p6NYhzFbc1oOJDVeBUCCVm39JJsMGqyIaOBqSum",
	"D7OCfWHdBiMXypmCkMnusAgD7MIoVlACW95K3+N7olRDCTukmBtObuvbWtD5gE4FCw3cpVAJQlh6c7Wb",
	"VMq01lUhk3yjDyz7yPIVghInMyPTfOSnyVQKzkQpgfBXh5E/efMiRdo7RpfqVQhbq1OHahMz7dk4cnpP",
	"uIQG3lUnZ1UHW1VnA/pEG2EM3sUs1oe6vqTwsdNSsuy2xZAs7/+Yc98EVXZ6wvNqNtIOqRjuS2C1itD8",
	"gMQWzCV0Yjv+tqzXN2ZxFdePKfvLAv4XEo52jtXeml7mrXCEg2/IT5OTxWcmeZdUi4wkjfdgZW5aACU4",
	"1co/Yc3BkpbQQMPUYXn2stf/6Hyc1bcy8WHdjLE1rtR9fX+p6nfPxdJ3zhR5qf5zRU+qiKyy0lNmFqWj",
	"LTSwZg5NVf3s6vpO88KkA95q34KQfui0H2TLHunLY+BsRD/v/wYAAP//D22wqn8GAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}