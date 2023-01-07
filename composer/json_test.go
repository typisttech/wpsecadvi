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
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockStringer struct {
	str string
}

func (ms mockStringer) String() string {
	return ms.str
}

func TestJSON_AddConflict(t *testing.T) {
	type fields struct {
		conflicts map[string]Link
	}
	type args struct {
		l Link
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "single",
			fields: fields{conflicts: nil},
			args: args{l: Link{
				name:        "foo",
				constraints: mockStringer{str: ">=1.2.3,<2.2.3"},
			}},
		},
		{
			name: "multiple",
			fields: fields{
				conflicts: map[string]Link{
					"bar": {
						name:        "bar",
						constraints: mockStringer{str: ">=9.8.7"},
					},
				},
			},
			args: args{l: Link{
				name:        "foo",
				constraints: mockStringer{str: ">=1.2.3,<2.2.3"},
			}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			j := &JSON{conflicts: tc.fields.conflicts}

			j.AddConflict(tc.args.l)

			assert.Contains(t, j.conflicts, tc.args.l.name)
			assert.Equal(t, tc.args.l, j.conflicts[tc.args.l.name])
		})
	}
}

func TestJSON_AddConflict_repetitive(t *testing.T) {
	foo := Link{
		name:        "foo",
		constraints: mockStringer{str: ">=1.2.3,<2.2.3"},
	}
	bar := Link{
		name:        "bar",
		constraints: mockStringer{str: ">=9.8.7"},
	}

	j := &JSON{}

	j.AddConflict(foo)
	j.AddConflict(bar)

	assert.Contains(t, j.conflicts, foo.name)
	assert.Equal(t, foo, j.conflicts[foo.name])
	assert.Contains(t, j.conflicts, bar.name)
	assert.Equal(t, bar, j.conflicts[bar.name])
}

func TestJSON_MarshalJSON(t *testing.T) {
	type fields struct {
		conflicts map[string]Link
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "single",
			fields: fields{conflicts: map[string]Link{
				"foo/bar": {
					name:        "foo/bar",
					constraints: mockStringer{str: ">=1,<2.2"},
				},
			}},
			want: `{"conflicts": {"foo/bar": ">=1,<2.2"}}`,
		},
		{
			name: "multiple",
			fields: fields{conflicts: map[string]Link{
				"foo/bar": {
					name:        "foo/bar",
					constraints: mockStringer{str: ">=1,<2.2"},
				},
				"bar/bar": {
					name:        "bar/bar",
					constraints: mockStringer{str: ">=1.2.3,<2.2.3|>=8.8.8|<=9.9.9"},
				},
				"baz/bar": {
					name:        "baz/bar",
					constraints: mockStringer{str: "*"},
				},
			}},
			want: `{"conflicts": {"baz/bar": "*", "foo/bar": ">=1,<2.2", "bar/bar": ">=1.2.3,<2.2.3|>=8.8.8|<=9.9.9"}}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			j := &JSON{
				conflicts: tc.fields.conflicts,
			}

			got, err := json.Marshal(j)

			if err != nil {
				t.Errorf("Unexpected error %v", err)
				return
			}

			assert.JSONEq(t, tc.want, string(got))
		})
	}
}

func TestJSON_Merge(t *testing.T) {
	type fields struct {
		conflicts map[string]Link
	}
	type args struct {
		other string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "single merge empty",
			fields: fields{conflicts: map[string]Link{
				"foo/bar": {
					name:        "foo/bar",
					constraints: mockStringer{str: ">=1,<2.2"},
				},
			}},
			args:    args{other: "{}"},
			want:    `{"conflicts": {"foo/bar": ">=1,<2.2"}}`,
			wantErr: false,
		},
		{
			name: "multiple merge empty",
			fields: fields{conflicts: map[string]Link{
				"foo/bar": {
					name:        "foo/bar",
					constraints: mockStringer{str: ">=1,<2.2"},
				},
				"bar/bar": {
					name:        "bar/bar",
					constraints: mockStringer{str: ">=1.2.3,<2.2.3|>=8.8.8|<=9.9.9"},
				},
				"baz/bar": {
					name:        "baz/bar",
					constraints: mockStringer{str: "*"},
				},
			}},
			args:    args{other: "{}"},
			want:    `{"conflicts": {"baz/bar": "*", "foo/bar": ">=1,<2.2", "bar/bar": ">=1.2.3,<2.2.3|>=8.8.8|<=9.9.9"}}`,
			wantErr: false,
		},
		{
			name: "merge",
			fields: fields{conflicts: map[string]Link{
				"foo/bar": {
					name:        "foo/bar",
					constraints: mockStringer{str: ">=1,<2.2"},
				},
				"bar/bar": {
					name:        "bar/bar",
					constraints: mockStringer{str: ">=1.2.3,<2.2.3|>=8.8.8|<=9.9.9"},
				},
				"baz/bar": {
					name:        "baz/bar",
					constraints: mockStringer{str: "*"},
				},
			}},
			args:    args{other: `{"name": "foo/bar"}`},
			want:    `{"name": "foo/bar", "conflicts": {"baz/bar": "*", "foo/bar": ">=1,<2.2", "bar/bar": ">=1.2.3,<2.2.3|>=8.8.8|<=9.9.9"}}`,
			wantErr: false,
		},
		{
			name: "override conflict",
			fields: fields{conflicts: map[string]Link{
				"foo/bar": {
					name:        "foo/bar",
					constraints: mockStringer{str: ">=1,<2.2"},
				},
				"bar/bar": {
					name:        "bar/bar",
					constraints: mockStringer{str: ">=1.2.3,<2.2.3|>=8.8.8|<=9.9.9"},
				},
				"baz/bar": {
					name:        "baz/bar",
					constraints: mockStringer{str: "*"},
				},
			}},
			args:    args{other: `{"name": "foo/bar", "conflicts": {"baz/bar": "1.2.3"}}`},
			want:    `{"name": "foo/bar", "conflicts": {"baz/bar": "*", "foo/bar": ">=1,<2.2", "bar/bar": ">=1.2.3,<2.2.3|>=8.8.8|<=9.9.9"}}`,
			wantErr: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			j := JSON{
				conflicts: tc.fields.conflicts,
			}
			got, err := j.Merge([]byte(tc.args.other))

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
