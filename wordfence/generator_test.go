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
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/typisttech/wpsecadvi/composer"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type subSearcher struct {
}

func (ss subSearcher) Search(t composer.PackageType, slug string) []string {
	return []string{
		"sub-" + string(t) + "/" + slug,
		"sub2-" + string(t) + "/" + slug,
	}
}

func TestGenerator_Generate(t *testing.T) {
	subSearcher := subSearcher{}

	type args struct {
		ignores []string
	}
	type fields struct {
		searcher composer.Searcher
	}
	tests := []struct {
		name    string
		fixture string
		fields  fields
		args    args
		want    string
	}{
		{
			name:    "production_ignore_id",
			fixture: "testdata/production.json",
			fields:  fields{searcher: subSearcher},
			args:    args{ignores: []string{"823e4567-e89b-12d3-a456-426655440000"}},
			want:    "testdata/production.composer.json.golden",
		},
		{
			name:    "production_ignore_cve",
			fixture: "testdata/production.json",
			fields:  fields{searcher: subSearcher},
			args:    args{ignores: []string{"CVE-2022-8888"}},
			want:    "testdata/production.composer.json.golden",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				file, _ := os.ReadFile(tc.fixture)
				fmt.Fprint(w, string(file))
			}))

			g := Generator{
				client: Client{
					httpClient: svr.Client(),
					url:        svr.URL,
				},
				searcher: tc.fields.searcher,
			}
			got, err := g.Generate(tc.args.ignores)
			if !assert.NoError(t, err) {
				return
			}

			gotJSON, err := json.Marshal(got)

			if !assert.NoError(t, err) {
				return
			}

			golden, err := os.ReadFile(tc.want)
			if !assert.NoError(t, err) {
				return
			}

			assert.JSONEq(t, string(golden), string(gotJSON))
		})
	}
}
