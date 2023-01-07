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

package composer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPrefixedSearcher(t *testing.T) {
	type args struct {
		t      PackageType
		prefix string
	}
	tests := []struct {
		name string
		args args
		want PrefixedSearcher
	}{
		{
			name: "happy_path",
			args: args{t: WPPlugin, prefix: "my-prefix"},
			want: PrefixedSearcher{packageType: WPPlugin, prefix: "my-prefix"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equalf(
				t,
				tc.want,
				NewPrefixedSearcher(tc.args.t, tc.args.prefix),
				"NewPrefixedSearcher(%v, %v)",
				tc.args.t,
				tc.args.prefix,
			)
		})
	}
}

func TestPrefixedSearcher_Search(t *testing.T) {
	type fields struct {
		packageType PackageType
		prefix      string
	}
	type args struct {
		t    PackageType
		slug string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name:   "happy_path",
			fields: fields{packageType: WPPlugin, prefix: "my-prefix"},
			args:   args{t: WPPlugin, slug: "foo-bar"},
			want:   []string{"my-prefix/foo-bar"},
		},
		{
			name:   "bad_package_type",
			fields: fields{packageType: WPPlugin, prefix: "my-prefix"},
			args:   args{t: WPTheme, slug: "foo-bar"},
			want:   nil, // Not expecting errors.
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ps := PrefixedSearcher{
				packageType: tc.fields.packageType,
				prefix:      tc.fields.prefix,
			}
			assert.Equalf(t, tc.want, ps.Search(tc.args.t, tc.args.slug), "Search(%v, %v)", tc.args.t, tc.args.slug)
		})
	}
}
