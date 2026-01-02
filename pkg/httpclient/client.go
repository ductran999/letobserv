package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type ReqHeader = map[string]string

// Client defines a simple HTTP client interface for making
// RESTful requests with optional headers and JSON payloads.
type Client interface {
	// Get sends an HTTP GET request to the specified URL with optional headers.
	Get(ctx context.Context, url string, headers ReqHeader) (*http.Response, error)

	// Post sends an HTTP POST request with a JSON-serializable body and optional headers.
	Post(ctx context.Context, url string, reqBody any, headers ReqHeader) (*http.Response, error)

	// Put sends an HTTP PUT request with a JSON-serializable body and optional headers.
	Put(ctx context.Context, url string, reqBody any, headers ReqHeader) (*http.Response, error)

	// Put sends an HTTP PUT request with a JSON-serializable body and optional headers.
	Patch(ctx context.Context, url string, reqBody any, headers ReqHeader) (*http.Response, error)

	// Delete sends an HTTP DELETE request to the specified URL with optional JSON body and headers.
	Delete(ctx context.Context, url string, reqBody any, headers ReqHeader) (*http.Response, error)
}

type httpClient struct {
	client *http.Client
}

func New() Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	clientTransport := otelhttp.NewTransport(
		transport,
		otelhttp.WithSpanNameFormatter(func(operation string, req *http.Request) string {
			fmt.Println(operation)
			urlPath := req.URL.Host
			return fmt.Sprintf("%s %s", req.Method, urlPath)
		}),
	)

	return &httpClient{
		client: &http.Client{
			Transport: clientTransport,
		},
	}
}

func (h *httpClient) Get(ctx context.Context, url string, headers ReqHeader) (*http.Response, error) {
	return h.doRequest(ctx, http.MethodGet, url, nil, headers)
}

func (h *httpClient) Post(ctx context.Context, url string, reqBody any, headers ReqHeader) (*http.Response, error) {
	return h.doRequest(ctx, http.MethodPost, url, reqBody, headers)
}

func (h *httpClient) Put(ctx context.Context, url string, reqBody any, headers ReqHeader) (*http.Response, error) {
	return h.doRequest(ctx, http.MethodPut, url, reqBody, headers)
}

func (h *httpClient) Patch(ctx context.Context, url string, reqBody any, headers ReqHeader) (*http.Response, error) {
	return h.doRequest(ctx, http.MethodPatch, url, reqBody, headers)
}

func (h *httpClient) Delete(ctx context.Context, url string, reqBody any, headers ReqHeader) (*http.Response, error) {
	return h.doRequest(ctx, http.MethodDelete, url, reqBody, headers)
}

// doRequest is the core method that executes an HTTP request.
func (h *httpClient) doRequest(
	ctx context.Context,
	method string,
	url string,
	reqBody any,
	headers ReqHeader,
) (*http.Response, error) {
	var body io.Reader

	// Marshal reqBody to JSON if not nil
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidBody, err)
		}
		body = bytes.NewReader(b)
	}

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateRequest, err)
	}

	// Set default Content-Type for JSON requests
	if headers == nil {
		headers = make(map[string]string)
	}
	if _, ok := headers["Content-Type"]; !ok && reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Apply per-request headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Perform the request
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDoRequest, err)
	}

	return resp, nil
}
