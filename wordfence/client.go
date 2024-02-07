package wordfence

import (
	"net/http"
	"net/url"
)

const (
	// ProductionFeed is the data feed with detailed records that have been fully
	// analyzed by the Wordfence team.
	ProductionFeed = "https://www.wordfence.com/api/intelligence/v2/vulnerabilities/production"

	// ScannerFeed is the data feed with minimal format that provides detection
	// information for newly discovered vulnerabilities that are actively being
	// researched in addition to those included in the ProductionFeed.
	ScannerFeed = "https://www.wordfence.com/api/intelligence/v2/vulnerabilities/scanner"
)

type Client struct {
	Client *http.Client
	URL    *url.URL
}

type ClientOpt func(c *Client)

func NewClient(opts ...ClientOpt) *Client {
	u, _ := url.Parse(ProductionFeed)
	client := &Client{
		Client: &http.Client{},
		URL:    u,
	}
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// NewProductionFeedClient returns a new Wordfence API client for fetching the
// ProductionFeed. If a nil httpClient is provided, a new [http.Client] will be
// used.
func NewProductionFeedClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	u, _ := url.Parse(ProductionFeed)

	return &Client{
		Client: httpClient,
		URL:    u,
	}
}

// NewScannerFeedClient returns a new Wordfence API client for fetching the
// ScannerFeed. If a nil httpClient is provided, a new [http.Client] will be
// used.
func NewScannerFeedClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	u, _ := url.Parse(ScannerFeed)

	return &Client{
		Client: httpClient,
		URL:    u,
	}
}
