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

package composer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_jsonMerge(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "empty json",
			args:    args{a: `{"foo": "bar"}`, b: `{}`},
			want:    `{"foo": "bar"}`,
			wantErr: false,
		},
		{
			name:    "normal",
			args:    args{a: `{"foo": "bar"}`, b: `{"baz": "qax"}`},
			want:    `{"foo": "bar", "baz":"qax"}`,
			wantErr: false,
		},
		{
			name:    "merge",
			args:    args{a: `{"foo": "bar", "baz": "qax"}`, b: `{"foo": "xxx", "baz": "yyy"}`},
			want:    `{"foo": "bar", "baz":"qax"}`,
			wantErr: false,
		},
		{
			name:    "override array",
			args:    args{a: `{"foo": ["aaa", "bbb"]}`, b: `{"foo": ["ccc", "ddd"]}`},
			want:    `{"foo": ["aaa", "bbb"]}`,
			wantErr: false,
		},
		{
			name:    "invalid json a",
			args:    args{a: `{`, b: `{"foo": "bar"}`},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid json b",
			args:    args{a: `{"foo": "bar"}`, b: `}`},
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty a",
			args:    args{a: "", b: `{"foo": "bar"}`},
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty b",
			args:    args{a: `{"foo": "bar"}`, b: ""},
			want:    "",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := jsonMerge([]byte(tc.args.a), []byte(tc.args.b))

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			if !assert.NoError(t, err) {
				return
			}

			assert.JSONEq(t, tc.want, string(got))
		})
	}
}
