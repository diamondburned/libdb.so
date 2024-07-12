package httpfs

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

type httpClient struct {
	client *http.Client
	url    *url.URL
}

// Timeout is the default timeout for http requests.
const Timeout = 8 * time.Second

func (c *httpClient) get(filepath string, _ FileInfo) (io.ReadCloser, error) {
	u := *c.url
	u.Path = path.Join(u.Path, filepath)

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
