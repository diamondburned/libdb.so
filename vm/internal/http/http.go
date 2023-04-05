package http

import (
	"context"
	"io"
	"net/http"
)

// Method describes an HTTP method.
type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
	PATCH  Method = "PATCH"
)

// Header returns the value of the header with the given key.
type Header = http.Header

// Request represents an HTTP request.
type Request struct {
	Method Method
	URL    string
	Header Header
	Body   io.Reader
}

// Response represents an HTTP response.
type Response struct {
	Status int
	Header Header
	Body   io.ReadCloser
}

// IsOK returns true if the response status is 2xx.
func (r *Response) IsOK() bool {
	return r.Status >= 200 && r.Status < 300
}

// RequestDoer is an interface for making HTTP requests.
type RequestDoer interface {
	Do(ctx context.Context, req *Request) (*Response, error)
}

// Client is a simple HTTP client.
type Client struct {
	Doer RequestDoer
}

// DefaultClient is the default HTTP client. It is set during init() to a client
// for the native platform.
var DefaultClient = &Client{}

// NewClient returns a new Client.
func NewClient(doer RequestDoer) *Client {
	return &Client{Doer: doer}
}

// Do makes an HTTP request.
func (c *Client) Do(ctx context.Context, req *Request) (*Response, error) {
	return c.Doer.Do(ctx, req)
}

// Get makes an HTTP GET request.
func (c *Client) Get(ctx context.Context, url string, header map[string][]string) (*Response, error) {
	return c.Do(ctx, &Request{
		Method: GET,
		URL:    url,
		Header: header,
	})
}
