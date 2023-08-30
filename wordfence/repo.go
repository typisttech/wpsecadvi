/*
 * Copyright (c) 2023 Typist Tech Limited
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
	"fmt"
	"slices"

	"github.com/typisttech/wpsecadvi/wp"
)

// excludeFunc returns true if a [Vulnerability] should be excluded from [Repo.Get] result.
type excludeFunc func(v Vulnerability) bool

type fetcher interface {
	fetch(ctx context.Context) (map[string]Vulnerability, error)
}

type Repo struct {
	Client       fetcher
	discardFuncs []excludeFunc
}

// WhereIDNotIn excludes [Vulnerability] from [Repo.Get] results by IDs.
func (r *Repo) WhereIDNotIn(ids ...string) *Repo {
	return r.withFilterFn(func(v Vulnerability) bool {
		return slices.Contains(ids, v.ID)
	})
}

// WhereCVENotIn excludes [Vulnerability] from [Repo.Get] results by CVEs.
func (r *Repo) WhereCVENotIn(cves ...string) *Repo {
	return r.withFilterFn(func(v Vulnerability) bool {
		return slices.Contains(cves, v.CVE)
	})
}

func (r *Repo) withFilterFn(fn excludeFunc) *Repo {
	r.discardFuncs = append(r.discardFuncs, fn)
	return r
}

func (r *Repo) Get(ctx context.Context) ([]*wp.Entity, error) {
	vsMap, err := r.Client.fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vulnerabilities form WordFence: %w", err)
	}

	vs := make(vulnerabilities, 0, len(vsMap))
	for _, v := range vsMap {
		vs = append(vs, v)
	}

	for _, fn := range r.discardFuncs {
		vs = slices.DeleteFunc(vs, fn)
	}

	return vs.entities()
}