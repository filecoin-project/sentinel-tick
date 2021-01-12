package quotetracker

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type quoteResponse interface {
	Quote() (Quote, error)
}

var globalHTTPClient = &http.Client{}

func request(ctx context.Context, url string, query url.Values, headers http.Header, response quoteResponse) (Quote, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		url,
		nil,
	)
	if err != nil {
		return Quote{}, err
	}

	if headers != nil {
		req.Header = headers
	}
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	resp, err := globalHTTPClient.Do(req)
	if err != nil {
		return Quote{}, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Quote{}, err
	}
	err = json.Unmarshal(respBody, response)
	if err != nil {
		return Quote{}, err
	}
	return response.Quote()
}
