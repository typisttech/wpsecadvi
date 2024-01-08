/*
 * Copyright (c) 2022-2023 Typist Tech Limited
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package wordfence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"golang.org/x/net/context/ctxhttp"
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

// excludeFunc returns true if a [Vulnerability] should be excluded from [Repo.Get] result.
type excludeFunc func(v Vulnerability) bool

type Client struct {
	// HTTPClient provides a http.Client fetch Wordfence data feed.
	// If the client is nil, http.DefaultClient is used.
	HTTPClient *http.Client
	// URL to the Wordfence data feed.
	// If the URL is empty, ProductionFeed is used.
	URL string

	excludeFuncs []excludeFunc
}

// WhereIDNotIn excludes [Vulnerability] results by IDs.
func (c *Client) WhereIDNotIn(ids ...string) *Client {
	return c.withExcludeFunc(func(v Vulnerability) bool {
		return slices.Contains(ids, v.ID)
	})
}

// WhereCVENotIn excludes [Vulnerability] by CVEs.
func (c *Client) WhereCVENotIn(cves ...string) *Client {
	return c.withExcludeFunc(func(v Vulnerability) bool {
		return slices.Contains(cves, v.CVE)
	})
}

func (c *Client) withExcludeFunc(fn excludeFunc) *Client {
	c.excludeFuncs = append(c.excludeFuncs, fn)
	return c
}

func (c *Client) fetch(ctx context.Context) (vulnerabilities, error) {
	url := c.URL
	if url == "" {
		url = ProductionFeed
	}

	vsMap, err := get[map[string]Vulnerability](ctx, c.HTTPClient, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vulnerabilities from Wordfence feed %s: %w", url, err)
	}

	vs := make(vulnerabilities, 0, len(vsMap))
	for _, v := range vsMap {
		vs = append(vs, v)
	}

	for _, fn := range c.excludeFuncs {
		vs = slices.DeleteFunc(vs, fn)
	}

	if len(vs) == 0 {
		return nil, fmt.Errorf("no vulnerabilities found from Wordfence feed %s", url)
	}

	return vs, nil
}

func get[T any](ctx context.Context, client *http.Client, url string) (T, error) {
	var out T

	res, err := ctxhttp.Get(ctx, client, url)
	if err != nil {
		return out, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return out, fmt.Errorf("HTTP GET request failed with %s", res.Status)
	}

	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return out, err
	}

	return out, nil
}