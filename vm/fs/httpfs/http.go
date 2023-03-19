package httpfs

import (
	"io"
	"net/http"
)

type httpClient struct {
	client http.Client
}

func (c *httpClient) get(path string, info FileInfo) (io.ReadCloser, error) {
	resp, err := c.client.Get(path)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
