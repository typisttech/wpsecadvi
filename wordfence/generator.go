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
	"github.com/typisttech/wpsecadvi/composer"
	"github.com/typisttech/wpsecadvi/semver"
)

type Generator struct {
	client   Client
	searcher composer.Searcher
}

func NewGenerator(client Client, searcher composer.Searcher) Generator {
	return Generator{
		client:   client,
		searcher: searcher,
	}
}

func (g Generator) Generate(ignores []string) (composer.JSON, error) {
	json := composer.JSON{}

	vulns, err := g.client.fetch()
	if err != nil {
		return json, err
	}

	links := make([]composer.Link, 0)
	for t, ranges := range vulns.softwareRanges(ignores) {
		for slug, rs := range ranges {
			// Just skip WPMU.
			if t == composer.WPCore && slug == "wpmu" {
				continue
			}

			ns := g.searcher.Search(t, string(slug))
			if len(ns) == 0 {
				continue
			}

			cs := semver.Constraints{}
			for _, r := range rs {
				cs.Or(r)
			}

			for _, n := range ns {
				if n == "" {
					continue
				}

				l := composer.NewLink(n, cs)

				links = append(links, l)
			}
		}
	}

	for _, l := range links {
		json.AddConflict(l)
	}

	return json, nil
}
