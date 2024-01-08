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
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestClient_fetch(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		want    vulnerabilities
		wantErr bool
	}{
		{
			"production example",
			"testdata/production.example.json",
			productionExample,
			false,
		},
		{
			"production multiple",
			"testdata/production.multiple.json",
			productionMultiple,
			false,
		},
		{
			"scanner example",
			"testdata/scanner.example.json",
			scannerExample,
			false,
		},
		{
			"empty",
			"testdata/empty.json",
			vulnerabilities{},
			false,
		},
		{
			"empty",
			"testdata/not-json.json",
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				file, _ := os.ReadFile(tt.fixture)
				_, _ = w.Write(file)
			}))
			defer ts.Close()

			c := Client{
				HTTPClient: ts.Client(),
				URL:        ts.URL,
			}
			ctx := context.Background()

			got, err := c.fetch(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			diff := cmp.Diff(
				tt.want,
				got,
				cmpopts.SortSlices(func(a, b Vulnerability) bool {
					return a.ID < b.ID
				}),
			)
			if diff != "" {
				t.Errorf("fetch() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestClient_fetch_cancelable(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		file, _ := os.ReadFile("testdata/production.example.json")
		_, _ = w.Write(file)

		time.Sleep(1 * time.Second)
	}))
	defer ts.Close()

	c := Client{
		HTTPClient: ts.Client(),
		URL:        ts.URL,
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	_, err := c.fetch(ctx)

	if !errors.Is(err, context.Canceled) {
		t.Errorf("fetch() error = %v, want %v", err, context.Canceled)
	}
}

func TestClient_fetch_http_error(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		file, _ := os.ReadFile("testdata/production.example.json")
		_, _ = w.Write(file)
	}))
	defer ts.Close()

	c := Client{
		HTTPClient: ts.Client(),
		URL:        ts.URL,
	}
	ctx := context.Background()

	_, err := c.fetch(ctx)

	if err == nil {
		t.Errorf("fetch() error = %v, wantErr true", err)
	}
}

func TestClient_WhereIDNotIn(t *testing.T) {
	tests := []struct {
		fixture string
		ids     []string
		want    vulnerabilities
		wantErr bool
	}{
		{
			fixture: "testdata/production.multiple.json",
			ids:     nil,
			want:    productionMultiple,
			wantErr: false,
		},
		{
			fixture: "testdata/production.multiple.json",
			ids:     []string{},
			want:    productionMultiple,
			wantErr: false,
		},
		{
			fixture: "testdata/production.multiple.json",
			ids:     []string{"014da588-9494-493e-8659-590b8e8c14a6"}, // wpgsi & wpgsiProfessional
			want:    append(productionMultiple[0:2], productionMultiple[3:]...),
			wantErr: false,
		},
		//{
		//	fixture: "testdata/production.multiple.json",
		//	ids: []string{
		//		"0114f098-713d-4eef-8643-901f607375de", // core
		//		"014da588-9494-493e-8659-590b8e8c14a6", // wpgsi & wpgsiProfessional
		//	},
		//	want:    append(productionMultiple[1:2], productionMultiple[4:]...),
		//	wantErr: false,
		//},
		//{
		//	fixture: "testdata/production.multiple.json",
		//	ids: []string{
		//		"0114f098-713d-4eef-8643-901f607375de", // core
		//		"01179ac2-ad68-4a5d-af67-70d57ed611d2", // simpleShippingEdd
		//		"014da588-9494-493e-8659-590b8e8c14a6", // wpgsi & wpgsiProfessional
		//		"06fee60a-e96c-49ce-9007-0d402ef46d72", // dtChocolate
		//	},
		//	want:    nil,
		//	wantErr: true,
		//},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				file, _ := os.ReadFile(tt.fixture)
				_, _ = w.Write(file)
			}))
			defer ts.Close()

			c := Client{
				HTTPClient: ts.Client(),
				URL:        ts.URL,
			}
			ctx := context.Background()

			got, err := c.WhereIDNotIn(tt.ids...).fetch(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("WhereIDNotIn().fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			diff := cmp.Diff(
				tt.want,
				got,
				cmpopts.SortSlices(func(a, b Vulnerability) bool {
					return a.ID < b.ID
				}),
				cmpopts.SortSlices(func(a, b Software) bool {
					return a.Type+a.Slug < b.Type+b.Slug
				}),
			)
			if diff != "" {
				t.Errorf("WhereIDNotIn().fetch() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}