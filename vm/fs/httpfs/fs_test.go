package httpfs

import (
	"fmt"
	"net/url"
	"path"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestReconcileURL(t *testing.T) {
	const appendPath = "/a"

	tests := []struct {
		baseURL  string
		fsURL    string
		expected string
		appended string
	}{
		{
			baseURL:  ".",
			fsURL:    "https://libdb.so/docs/_docsfs.json",
			expected: "https://libdb.so/docs",
			appended: "https://libdb.so/docs/a",
		},
		{
			baseURL:  "https://docs.libdb.so/",
			fsURL:    "https://docs.libdb.so/_docsfs.json",
			expected: "https://docs.libdb.so/",
			appended: "https://docs.libdb.so/a",
		},
		{
			baseURL:  "https://docs.libdb.so/.",
			fsURL:    "https://docs.libdb.so/_docsfs.json",
			expected: "https://docs.libdb.so/.",
			appended: "https://docs.libdb.so/a",
		},
		{
			baseURL:  "https://docs.libdb.so",
			fsURL:    "https://docs.libdb.so/_docsfs.json",
			expected: "https://docs.libdb.so/",
			appended: "https://docs.libdb.so/a",
		},
		{
			baseURL:  "https://docs.libdb.so",
			fsURL:    "https://libdb.so/docs/_docsfs.json",
			expected: "https://docs.libdb.so/docs",
			appended: "https://docs.libdb.so/docs/a",
		},
		{
			baseURL:  "https://docs.libdb.so/",
			fsURL:    "https://libdb.so/docs/_docsfs.json",
			expected: "https://docs.libdb.so/",
			appended: "https://docs.libdb.so/a",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint("case-", i+1), func(t *testing.T) {
			baseURL := mustURL(test.baseURL)
			fsURL := mustURL(test.fsURL)

			got := reconcileURL(*baseURL, *fsURL)
			assert.Equal(t, test.expected, got.String(), "expected")

			got.Path = path.Join(got.Path, appendPath)
			assert.Equal(t, test.appended, got.String(), "appended")
		})
	}
}

func mustURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
