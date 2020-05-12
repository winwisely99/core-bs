package fetcher

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/common/logger"
	rh "github.com/hashicorp/go-retryablehttp"
	"net/http"
)

type fetcherClient struct {
	*rh.Client
}

// NewClient returns new retryable http client
func NewClient(l *logger.Logger) *fetcherClient {
	cl := rh.NewClient()
	cl.Logger = l
	return &fetcherClient{cl}
}

// Fetch making http request given url, method and request body.
func (f *fetcherClient) Fetch(url, method string, body interface{}, headers map[string]string) (*http.Response, error) {
	req, err := rh.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if headers == nil {
		headers = map[string]string{
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"Accept-Encoding": "gzip,deflate,br",
		}
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return f.Do(req)
}
