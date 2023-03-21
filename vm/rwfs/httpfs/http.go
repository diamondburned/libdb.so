package httpfs

import (
	"io"
	"net/http"
	"path"
)

type httpClient struct {
	client   http.Client
	basePath string
}

func (c *httpClient) get(filepath string, info FileInfo) (io.ReadCloser, error) {
	filepath = path.Join(c.basePath, filepath)

	resp, err := c.client.Get(filepath)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
