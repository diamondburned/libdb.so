package httpfs

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"hash"
	"io"
	"net/http"
)

func wrapReaderForCache(path string, info FileInfo, r io.ReadCloser) io.ReadCloser {
	r = &checksumReadCloser{r: r, hash: info.Hash, hasher: sha256.New()}
	return r
}

type checksumReadCloser struct {
	r      io.ReadCloser
	hash   []byte
	hasher hash.Hash
}

func (c *checksumReadCloser) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.hasher.Write(p[:n])

	if errors.Is(err, io.EOF) {
		gotHash := c.hasher.Sum(nil)
		if !bytes.Equal(gotHash, c.hash) {
			return n, errors.New("httpfs: checksum mismatch")
		}
	}

	return n, err
}

func (c *checksumReadCloser) Close() error {
	return c.r.Close()
}

type httpClient struct {
	client http.Client
}

func (c *httpClient) get(path string, info FileInfo) (io.ReadCloser, error) {
	resp, err := c.client.Get(path)
	if err != nil {
		return nil, err
	}

	r := wrapReaderForCache(path, info, resp.Body)
	return r, nil
}
