package httpfs

import (
	"context"
	"io"
	"net/http"
	"path"
	"time"
)

type httpClient struct {
	client   http.Client
	basePath string
}

// Timeout is the default timeout for http requests.
const Timeout = 8 * time.Second

func (c *httpClient) get(filepath string, info FileInfo) (io.ReadCloser, error) {
	filepath = path.Join(c.basePath, filepath)

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, filepath, nil)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
