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
	"reflect"
	"testing"
	"time"
)

func TestClient_fetch(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		want    map[string]Vulnerability
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
			map[string]Vulnerability{},
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
				w.Write(file)
			}))
			defer ts.Close()

			c := Client{
				HTTPClient: ts.Client(),
				URL:        ts.URL,
			}
			ctx := context.Background()

			got, err := c.fetch(ctx)

			if tt.wantErr && (err == nil) {
				t.Errorf("fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (err != nil) {
				t.Errorf("fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_fetch_cancelable(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		file, _ := os.ReadFile("testdata/production.example.json")
		w.Write(file)

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
		w.Write(file)
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