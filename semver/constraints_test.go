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

package semver

import (
	"github.com/Masterminds/semver/v3"
	"reflect"
	"testing"
)

var (
	wildcardRange = &Range{
		fromWildcard: true,
		toWildcard:   true,
	}

	v123, _ = semver.NewVersion("1.2.3")
	v133, _ = semver.NewVersion("1.3.3")
	v223, _ = semver.NewVersion("2.2.3")
	v233, _ = semver.NewVersion("2.3.3")
	v323, _ = semver.NewVersion("3.2.3")
	v423, _ = semver.NewVersion("4.2.3")
	v523, _ = semver.NewVersion("5.2.3")
	v533, _ = semver.NewVersion("5.3.3")
	v623, _ = semver.NewVersion("6.2.3")

	v123ToV133                   = &Range{fromVersion: v123, toVersion: v133}
	v123ToV223                   = &Range{fromVersion: v123, toVersion: v223}
	v123ToV233                   = &Range{fromVersion: v123, toVersion: v233}
	v123ToV423                   = &Range{fromVersion: v123, toVersion: v423}
	v123ToV623                   = &Range{fromVersion: v123, toVersion: v623}
	v123InclusiveToV623Inclusive = &Range{fromVersion: v123, fromInclusive: true, toVersion: v623, toInclusive: true}
	v223ToV323                   = &Range{fromVersion: v223, toVersion: v323}
	v223InclusiveToV323Inclusive = &Range{fromVersion: v223, fromInclusive: true, toVersion: v323, toInclusive: true}
	v223ToV323Inclusive          = &Range{fromVersion: v223, toVersion: v323, toInclusive: true}
	v223InclusiveToV323          = &Range{fromVersion: v223, fromInclusive: true, toVersion: v323}
	v223ToV423                   = &Range{fromVersion: v223, toVersion: v423}
	v223ToV623                   = &Range{fromVersion: v223, toVersion: v623}
	v423ToV523                   = &Range{fromVersion: v423, toVersion: v523}
	v423ToV623                   = &Range{fromVersion: v423, toVersion: v623}
	v323ToV523                   = &Range{fromVersion: v323, toVersion: v523}
	v533ToV623                   = &Range{fromVersion: v533, toVersion: v623}

	gtV223          = &Range{fromVersion: v223, toWildcard: true}
	gtV223Inclusive = &Range{fromVersion: v223, fromInclusive: true, toWildcard: true}
	gtV323          = &Range{fromVersion: v323, fromInclusive: false, toWildcard: true}
	gtV323Inclusive = &Range{fromVersion: v323, fromInclusive: true, toWildcard: true}
	gtV523          = &Range{fromVersion: v523, fromInclusive: false, toWildcard: true}

	ltV223          = &Range{toVersion: v223, fromWildcard: true}
	ltV223Inclusive = &Range{toVersion: v223, toInclusive: true, fromWildcard: true}
	ltV323          = &Range{toVersion: v323, toInclusive: false, fromWildcard: true}
	ltV323Inclusive = &Range{toVersion: v323, toInclusive: true, fromWildcard: true}
)

func TestConstraints_String(t *testing.T) {
	type fields struct {
		ranges []*Range
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "single",
			fields: fields{ranges: []*Range{v223ToV323Inclusive}},
			want:   ">2.2.3,<=3.2.3",
		},
		{
			name:   "single wildcard",
			fields: fields{ranges: []*Range{wildcardRange}},
			want:   "*",
		},
		{
			name:   "multiple",
			fields: fields{ranges: []*Range{v123ToV223, ltV223, gtV323Inclusive}},
			want:   "<2.2.3|>1.2.3,<2.2.3|>=3.2.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &Constraints{
				ranges: tt.fields.ranges,
			}
			if got := cs.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstraints_Or_wildcard(t *testing.T) {
	type fields struct {
		ranges []*Range
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Range
	}{
		{
			name:   "empty",
			fields: fields{ranges: []*Range{}},
			want:   []*Range{wildcardRange},
		},
		{
			name: "wildcard",
			fields: fields{ranges: []*Range{
				{fromWildcard: true, toWildcard: true},
			}},
			want: []*Range{wildcardRange},
		},
		{
			name:   "upper_bounded_non-inclusive",
			fields: fields{ranges: []*Range{ltV223}},
			want:   []*Range{wildcardRange},
		},
		{
			name:   "upper_bounded_inclusive",
			fields: fields{ranges: []*Range{ltV223Inclusive}},
			want:   []*Range{wildcardRange},
		},
		{
			name:   "lower_bounded_non-inclusive",
			fields: fields{ranges: []*Range{gtV223}},
			want:   []*Range{wildcardRange},
		},
		{
			name:   "lower_bounded_inclusive",
			fields: fields{ranges: []*Range{gtV223Inclusive}},
			want:   []*Range{wildcardRange},
		},
		{
			name: "reduce_to_single_wildcard",
			fields: fields{
				ranges: []*Range{
					v123ToV223,
					{
						fromVersion:   v323,
						fromInclusive: true,
						toVersion:     v423,
						toInclusive:   true,
					},
					{
						fromVersion: v523,
						toWildcard:  true,
					},
					{
						fromWildcard: true,
						toVersion:    v623,
					},
				},
			},
			want: []*Range{wildcardRange},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cs := &Constraints{
				ranges: tc.fields.ranges,
			}

			cs.Or(wildcardRange)

			if !reflect.DeepEqual(cs.ranges, tc.want) {
				t.Errorf("cs.ranges = %v, want %v", cs.ranges, tc.want)
			}
		})
	}
}

func TestConstraints_Or_upper_bounded(t *testing.T) {
	type fields struct {
		ranges []*Range
	}
	type args struct {
		r *Range
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Range
	}{
		{
			name:   "empty",
			fields: fields{ranges: []*Range{}},
			args:   args{r: ltV323Inclusive},
			want:   []*Range{ltV323Inclusive},
		},
		{
			name:   "one_unrelated",
			fields: fields{ranges: []*Range{v423ToV523}},
			args:   args{r: ltV323Inclusive},
			want:   []*Range{v423ToV523, ltV323Inclusive},
		},
		{
			name:   "two_unrelated",
			fields: fields{ranges: []*Range{v423ToV523, v533ToV623}},
			args:   args{r: ltV323},
			want:   []*Range{v423ToV523, v533ToV623, ltV323},
		},
		{
			name:   "wildcard",
			fields: fields{ranges: []*Range{wildcardRange}},
			args:   args{r: ltV323Inclusive},
			want:   []*Range{wildcardRange},
		},
		{
			name:   "same_inclusive",
			fields: fields{ranges: []*Range{ltV323Inclusive}},
			args:   args{r: ltV323Inclusive},
			want:   []*Range{ltV323Inclusive},
		},
		{
			name:   "same_non-inclusive",
			fields: fields{ranges: []*Range{ltV323}},
			args:   args{r: ltV323},
			want:   []*Range{ltV323},
		},
		{
			name:   "same_ranges_inclusive_r_non-inclusive",
			fields: fields{ranges: []*Range{ltV323Inclusive}},
			args:   args{r: ltV323},
			want:   []*Range{ltV323Inclusive},
		},
		{
			name:   "same ranges_non-inclusive_r_inclusive",
			fields: fields{ranges: []*Range{ltV323}},
			args:   args{r: ltV323Inclusive},
			want:   []*Range{ltV323Inclusive},
		},
		{
			name:   "superset",
			fields: fields{ranges: []*Range{ltV323}},
			args:   args{r: ltV223},
			want:   []*Range{ltV323},
		},
		{
			name:   "superset_inclusive",
			fields: fields{ranges: []*Range{ltV323Inclusive}},
			args:   args{r: ltV223Inclusive},
			want:   []*Range{ltV323Inclusive},
		},
		{
			name:   "superset_with_unrelated",
			fields: fields{ranges: []*Range{v423ToV523, ltV323Inclusive}},
			args:   args{r: ltV223Inclusive},
			want:   []*Range{v423ToV523, ltV323Inclusive},
		},
		{
			name:   "superset_with_overlap",
			fields: fields{ranges: []*Range{v123ToV423, ltV223Inclusive}},
			args:   args{r: ltV323Inclusive},
			want:   []*Range{v123ToV423, ltV323Inclusive},
		},
		{
			name:   "one_subset",
			fields: fields{ranges: []*Range{ltV223}},
			args:   args{r: ltV323},
			want:   []*Range{ltV323},
		},
		{
			name:   "one_subset_both_bounded",
			fields: fields{ranges: []*Range{v123ToV233}},
			args:   args{r: ltV323},
			want:   []*Range{ltV323},
		},
		{
			name:   "one_subset_with_unrelated",
			fields: fields{ranges: []*Range{v423ToV523, ltV223}},
			args:   args{r: ltV323},
			want:   []*Range{v423ToV523, ltV323},
		},
		{
			name:   "one_subset_with_overlap",
			fields: fields{ranges: []*Range{v123ToV423, ltV223}},
			args:   args{r: ltV323},
			want:   []*Range{v123ToV423, ltV323},
		},

		{
			name:   "lower_bounded_unrelated",
			fields: fields{ranges: []*Range{gtV323}},
			args:   args{r: ltV223},
			want:   []*Range{gtV323, ltV223},
		},
		{
			name:   "lower_bounded_overlapped",
			fields: fields{ranges: []*Range{gtV223}},
			args:   args{r: ltV323},
			want:   []*Range{gtV223, ltV323},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cs := &Constraints{
				ranges: tc.fields.ranges,
			}

			cs.Or(tc.args.r)

			if !reflect.DeepEqual(cs.ranges, tc.want) {
				t.Errorf("cs.ranges = %v, want %v", cs.ranges, tc.want)
			}
		})
	}
}

func TestConstraints_Or_lower_bounded(t *testing.T) {
	type fields struct {
		ranges []*Range
	}
	type args struct {
		r *Range
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Range
	}{
		{
			name:   "empty",
			fields: fields{ranges: []*Range{}},
			args:   args{r: gtV323Inclusive},
			want:   []*Range{gtV323Inclusive},
		},
		{
			name:   "one_unrelated",
			fields: fields{ranges: []*Range{v123ToV233}},
			args:   args{r: gtV323Inclusive},
			want:   []*Range{v123ToV233, gtV323Inclusive},
		},
		{
			name:   "two_unrelated",
			fields: fields{ranges: []*Range{v123ToV133, v123ToV223}},
			args:   args{r: gtV323},
			want:   []*Range{v123ToV133, v123ToV223, gtV323},
		},
		{
			name:   "wildcard",
			fields: fields{ranges: []*Range{wildcardRange}},
			args:   args{r: gtV323Inclusive},
			want:   []*Range{wildcardRange},
		},
		{
			name:   "same_inclusive",
			fields: fields{ranges: []*Range{gtV323Inclusive}},
			args:   args{r: gtV323Inclusive},
			want:   []*Range{gtV323Inclusive},
		},
		{
			name:   "same_non-inclusive",
			fields: fields{ranges: []*Range{gtV323}},
			args:   args{r: gtV323},
			want:   []*Range{gtV323},
		},
		{
			name:   "same_ranges_inclusive_r_non-inclusive",
			fields: fields{ranges: []*Range{gtV323Inclusive}},
			args:   args{r: gtV323},
			want:   []*Range{gtV323Inclusive},
		},
		{
			name:   "same ranges_non-inclusive_r_inclusive",
			fields: fields{ranges: []*Range{gtV323}},
			args:   args{r: gtV323Inclusive},
			want:   []*Range{gtV323Inclusive},
		},
		{
			name:   "superset",
			fields: fields{ranges: []*Range{gtV223}},
			args:   args{r: gtV323},
			want:   []*Range{gtV223},
		},
		{
			name:   "superset_inclusive",
			fields: fields{ranges: []*Range{gtV223Inclusive}},
			args:   args{r: gtV323Inclusive},
			want:   []*Range{gtV223Inclusive},
		},
		{
			name:   "superset_with_unrelated",
			fields: fields{ranges: []*Range{v123ToV133, gtV223Inclusive}},
			args:   args{r: gtV323Inclusive},
			want:   []*Range{v123ToV133, gtV223Inclusive},
		},
		{
			name:   "superset_with_overlap",
			fields: fields{ranges: []*Range{v123ToV423, gtV323Inclusive}},
			args:   args{r: gtV223Inclusive},
			want:   []*Range{v123ToV423, gtV223Inclusive},
		},
		{
			name:   "one_subset",
			fields: fields{ranges: []*Range{gtV323}},
			args:   args{r: gtV223},
			want:   []*Range{gtV223},
		},
		{
			name:   "one_subset_both_bounded",
			fields: fields{ranges: []*Range{v323ToV523}},
			args:   args{r: gtV223},
			want:   []*Range{gtV223},
		},
		{
			name:   "one_subset_with_unrelated",
			fields: fields{ranges: []*Range{v123ToV133, gtV323}},
			args:   args{r: gtV223},
			want:   []*Range{v123ToV133, gtV223},
		},
		{
			name:   "one_subset_with_overlap",
			fields: fields{ranges: []*Range{v123ToV423, gtV323}},
			args:   args{r: gtV223},
			want:   []*Range{v123ToV423, gtV223},
		},
		{
			name:   "upper_bounded_unrelated",
			fields: fields{ranges: []*Range{ltV223}},
			args:   args{r: gtV323},
			want:   []*Range{ltV223, gtV323},
		},
		{
			name:   "upper_bounded_overlapped",
			fields: fields{ranges: []*Range{ltV323}},
			args:   args{r: gtV223},
			want:   []*Range{ltV323, gtV223},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cs := &Constraints{
				ranges: tc.fields.ranges,
			}

			cs.Or(tc.args.r)

			if !reflect.DeepEqual(cs.ranges, tc.want) {
				t.Errorf("cs.ranges = %v, want %v", cs.ranges, tc.want)
			}
		})
	}
}

func TestConstraints_Or_bounded(t *testing.T) {
	type fields struct {
		ranges []*Range
	}
	type args struct {
		r *Range
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Range
	}{
		{
			name: "empty",
			fields: fields{
				ranges: []*Range{},
			},
			args: args{r: v223ToV323},
			want: []*Range{v223ToV323},
		},
		{
			name: "one_unrelated",
			fields: fields{
				ranges: []*Range{v123ToV133},
			},
			args: args{r: v223ToV323},
			want: []*Range{v123ToV133, v223ToV323},
		},
		{
			name: "two_unrelated",
			fields: fields{
				ranges: []*Range{v123ToV133, v423ToV523},
			},
			args: args{r: v223ToV323},
			want: []*Range{v123ToV133, v423ToV523, v223ToV323},
		},
		{
			name: "wildcard",
			fields: fields{
				ranges: []*Range{wildcardRange},
			},
			args: args{r: v223ToV323},
			want: []*Range{wildcardRange},
		},
		{
			name: "same_inclusive",
			fields: fields{
				ranges: []*Range{v223InclusiveToV323Inclusive},
			},
			args: args{r: v223InclusiveToV323Inclusive},
			want: []*Range{v223InclusiveToV323Inclusive},
		},
		{
			name: "same_non-inclusive",
			fields: fields{
				ranges: []*Range{v223ToV323},
			},
			args: args{r: v223ToV323},
			want: []*Range{v223ToV323},
		},
		{
			name: "same_left_non-inclusive_right_inclusive",
			fields: fields{
				ranges: []*Range{v223ToV323Inclusive},
			},
			args: args{r: v223ToV323Inclusive},
			want: []*Range{v223ToV323Inclusive},
		},
		{
			name: "same_left_inclusive_right_non-inclusive",
			fields: fields{
				ranges: []*Range{v223InclusiveToV323},
			},
			args: args{r: v223InclusiveToV323},
			want: []*Range{v223InclusiveToV323},
		},
		{
			name: "same_ranges_inclusive_r_non-inclusive",
			fields: fields{
				ranges: []*Range{v223InclusiveToV323Inclusive},
			},
			args: args{r: v223ToV323},
			want: []*Range{v223InclusiveToV323Inclusive},
		},
		{
			name: "same_ranges_non-inclusive_r_inclusive",
			fields: fields{
				ranges: []*Range{v223ToV323},
			},
			args: args{r: v223InclusiveToV323Inclusive},
			want: []*Range{v223InclusiveToV323Inclusive},
		},
		{
			name: "same_ranges_left_inclusive_right_inclusive_r_left_non-inclusive_right_inclusive",
			fields: fields{
				ranges: []*Range{v223InclusiveToV323Inclusive},
			},
			args: args{r: v223ToV323Inclusive},
			want: []*Range{v223InclusiveToV323Inclusive},
		},
		{
			name: "same_ranges_left_non-inclusive_right_inclusive_r_left_inclusive_right_inclusive",
			fields: fields{
				ranges: []*Range{v223ToV323Inclusive},
			},
			args: args{r: v223InclusiveToV323Inclusive},
			want: []*Range{v223InclusiveToV323Inclusive},
		},
		{
			name: "same_ranges_left_inclusive_right_non-inclusive_r_left_non-inclusive_right_non-inclusive",
			fields: fields{
				ranges: []*Range{v223InclusiveToV323},
			},
			args: args{r: v223ToV323},
			want: []*Range{v223InclusiveToV323},
		},
		{
			name: "same_ranges_left_non-inclusive_right_non-inclusive_r_left_inclusive_right_non-inclusive",
			fields: fields{
				ranges: []*Range{v223ToV323},
			},
			args: args{r: v223InclusiveToV323},
			want: []*Range{v223InclusiveToV323},
		},
		{
			name: "same_ranges_left_inclusive_right_inclusive_r_left_inclusive_right_non-inclusive",
			fields: fields{
				ranges: []*Range{v223InclusiveToV323Inclusive},
			},
			args: args{r: v223InclusiveToV323},
			want: []*Range{v223InclusiveToV323Inclusive},
		},
		{
			name: "same_ranges_left_inclusive_right_non-inclusive_r_left_inclusive_right_inclusive",
			fields: fields{
				ranges: []*Range{v223InclusiveToV323},
			},
			args: args{r: v223InclusiveToV323Inclusive},
			want: []*Range{v223InclusiveToV323Inclusive},
		},
		{
			name: "same_ranges_left_non-inclusive_right_inclusive_r_left_inclusive_right_non-inclusive",
			fields: fields{
				ranges: []*Range{v223ToV323Inclusive},
			},
			args: args{r: v223InclusiveToV323},
			want: []*Range{v223ToV323Inclusive, v223InclusiveToV323},
		},
		{
			name: "superset",
			fields: fields{
				ranges: []*Range{v123ToV623},
			},
			args: args{r: v323ToV523},
			want: []*Range{v123ToV623},
		},
		{
			name: "superset_inclusive",
			fields: fields{
				ranges: []*Range{v123InclusiveToV623Inclusive},
			},
			args: args{r: v223ToV623},
			want: []*Range{v123InclusiveToV623Inclusive},
		},
		{
			name: "superset_with_unrelated",
			fields: fields{
				ranges: []*Range{v123ToV623, gtV523},
			},
			args: args{r: v223ToV423},
			want: []*Range{v123ToV623, gtV523},
		},
		{
			name: "superset_with_overlap",
			fields: fields{
				ranges: []*Range{v123ToV623, gtV323},
			},
			args: args{r: v223ToV423},
			want: []*Range{v123ToV623, gtV323},
		},
		{
			name: "subset",
			fields: fields{
				ranges: []*Range{v223ToV423},
			},
			args: args{r: v123ToV623},
			want: []*Range{v123ToV623},
		},
		{
			name: "left overlap",
			fields: fields{
				ranges: []*Range{v323ToV523},
			},
			args: args{r: v223ToV423},
			want: []*Range{v323ToV523, v223ToV423},
		},
		{
			name: "right overlap",
			fields: fields{
				ranges: []*Range{v323ToV523},
			},
			args: args{r: v423ToV623},
			want: []*Range{v323ToV523, v423ToV623},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cs := &Constraints{
				ranges: tc.fields.ranges,
			}

			cs.Or(tc.args.r)

			if !reflect.DeepEqual(cs.ranges, tc.want) {
				t.Errorf("cs.ranges = %v, want %v", cs.ranges, tc.want)
			}
		})
	}
}

func TestNewWildcardRange(t *testing.T) {
	want := &Range{
		fromWildcard: true,
		toWildcard:   true,
	}

	got, err := NewWildcardRange()

	if err != nil {
		t.Errorf("NewWildcardRange() error = %v", err)
		return
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewWildcardRange() = %v, want %v", got, want)
	}

	if !got.fromWildcard {
		t.Errorf(
			"NewWildcardRange().fromWildcard got = %v, want %v",
			got.fromWildcard,
			true,
		)
	}

	if !got.toWildcard {
		t.Errorf(
			"NewWildcardRange().toWildcard got = %v, want %v",
			got.toWildcard,
			true,
		)
	}
}

func TestNewUpperBoundedRange(t *testing.T) {
	type args struct {
		version   string
		inclusive bool
	}
	type want struct {
		toString    string
		toInclusive bool
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name:    "inclusive",
			args:    args{version: "1.2.3", inclusive: true},
			want:    want{toString: "1.2.3", toInclusive: true},
			wantErr: false,
		},
		{
			name:    "non-inclusive",
			args:    args{version: "1.2.3", inclusive: false},
			want:    want{toString: "1.2.3", toInclusive: false},
			wantErr: false,
		},
		{
			name:    "major_minor",
			args:    args{version: "1.2", inclusive: true},
			want:    want{toString: "1.2.0", toInclusive: true},
			wantErr: false,
		},
		{
			name:    "major",
			args:    args{version: "1", inclusive: true},
			want:    want{toString: "1.0.0", toInclusive: true},
			wantErr: false,
		},
		{
			name:    "invalid_semantic_version",
			args:    args{version: "1.2.3.4", inclusive: true},
			want:    want{},
			wantErr: true,
		},
		{
			name:    "invalid_version",
			args:    args{version: "informational", inclusive: true},
			want:    want{},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewUpperBoundedRange(tc.args.version, tc.args.inclusive)

			if tc.wantErr {
				if err == nil {
					t.Errorf("NewUpperBoundedRange() error = %v, wantErr %v", err, tc.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewUpperBoundedRange() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if got.toInclusive != tc.want.toInclusive {
				t.Errorf(
					"NewUpperBoundedRange().toInclusive got = %v, want %v",
					got.toInclusive,
					tc.want.toInclusive,
				)
			}

			if got.toVersion.String() != tc.want.toString {
				t.Errorf(
					"NewUpperBoundedRange().toVersion.String() got = %v, want %v",
					got.toVersion.String(),
					tc.want.toString,
				)
			}

			if !got.fromWildcard {
				t.Errorf(
					"NewUpperBoundedRange().fromWildcard got = %v, want %v",
					got.fromWildcard,
					true,
				)
			}

			if got.toWildcard {
				t.Errorf(
					"NewUpperBoundedRange().toWildcard got = %v, want %v",
					got.toWildcard,
					false,
				)
			}
		})
	}
}

func TestNewLowerBoundedRange(t *testing.T) {
	type args struct {
		version   string
		inclusive bool
	}
	type want struct {
		fromString    string
		fromInclusive bool
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name:    "inclusive",
			args:    args{version: "1.2.3", inclusive: true},
			want:    want{fromString: "1.2.3", fromInclusive: true},
			wantErr: false,
		},
		{
			name:    "non-inclusive",
			args:    args{version: "1.2.3", inclusive: false},
			want:    want{fromString: "1.2.3", fromInclusive: false},
			wantErr: false,
		},
		{
			name:    "major_minor",
			args:    args{version: "1.2", inclusive: true},
			want:    want{fromString: "1.2.0", fromInclusive: true},
			wantErr: false,
		},
		{
			name:    "major",
			args:    args{version: "1", inclusive: true},
			want:    want{fromString: "1.0.0", fromInclusive: true},
			wantErr: false,
		},
		{
			name:    "invalid_semantic_version",
			args:    args{version: "1.2.3.4", inclusive: true},
			want:    want{},
			wantErr: true,
		},
		{
			name:    "invalid_version",
			args:    args{version: "informational", inclusive: true},
			want:    want{},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewLowerBoundedRange(tc.args.version, tc.args.inclusive)

			if tc.wantErr {
				if err == nil {
					t.Errorf("NewLowerBoundedRange() error = %v, wantErr %v", err, tc.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewLowerBoundedRange() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if got.fromInclusive != tc.want.fromInclusive {
				t.Errorf(
					"NewLowerBoundedRange().fromInclusive got = %v, want %v",
					got.fromInclusive,
					tc.want.fromInclusive,
				)
			}

			if got.fromVersion.String() != tc.want.fromString {
				t.Errorf(
					"NewLowerBoundedRange().fromVersion.String() got = %v, want %v",
					got.fromVersion.String(),
					tc.want.fromString,
				)
			}

			if got.fromWildcard {
				t.Errorf(
					"NewLowerBoundedRange().fromWildcard got = %v, want %v",
					got.fromWildcard,
					false,
				)
			}

			if !got.toWildcard {
				t.Errorf(
					"NewLowerBoundedRange().toWildcard got = %v, want %v",
					got.toWildcard,
					true,
				)
			}
		})
	}
}

func TestNewRange(t *testing.T) {
	type args struct {
		fromVersion   string
		fromInclusive bool
		toVersion     string
		toInclusive   bool
	}
	type want struct {
		fromString    string
		fromInclusive bool
		toString      string
		toInclusive   bool
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name:    "from_inclusive",
			args:    args{fromVersion: "1.2.3", fromInclusive: true, toVersion: "9.8.7", toInclusive: false},
			want:    want{fromString: "1.2.3", fromInclusive: true, toString: "9.8.7", toInclusive: false},
			wantErr: false,
		},
		{
			name:    "to_inclusive",
			args:    args{fromVersion: "1.2.3", fromInclusive: false, toVersion: "9.8.7", toInclusive: true},
			want:    want{fromString: "1.2.3", fromInclusive: false, toString: "9.8.7", toInclusive: true},
			wantErr: false,
		},
		{
			name:    "both_inclusive",
			args:    args{fromVersion: "1.2.3", fromInclusive: true, toVersion: "9.8.7", toInclusive: true},
			want:    want{fromString: "1.2.3", fromInclusive: true, toString: "9.8.7", toInclusive: true},
			wantErr: false,
		},
		{
			name:    "both_non-inclusive",
			args:    args{fromVersion: "1.2.3", fromInclusive: false, toVersion: "9.8.7", toInclusive: false},
			want:    want{fromString: "1.2.3", fromInclusive: false, toString: "9.8.7", toInclusive: false},
			wantErr: false,
		},
		{
			name:    "from_major_minor",
			args:    args{fromVersion: "1.2", fromInclusive: false, toVersion: "9.8.7", toInclusive: false},
			want:    want{fromString: "1.2.0", fromInclusive: false, toString: "9.8.7", toInclusive: false},
			wantErr: false,
		},
		{
			name:    "to_major_minor",
			args:    args{fromVersion: "1.2.3", fromInclusive: false, toVersion: "9.8", toInclusive: false},
			want:    want{fromString: "1.2.3", fromInclusive: false, toString: "9.8.0", toInclusive: false},
			wantErr: false,
		},
		{
			name:    "both_major_minor",
			args:    args{fromVersion: "1.2", fromInclusive: false, toVersion: "9.8", toInclusive: false},
			want:    want{fromString: "1.2.0", fromInclusive: false, toString: "9.8.0", toInclusive: false},
			wantErr: false,
		},
		{
			name:    "from_major",
			args:    args{fromVersion: "1", fromInclusive: false, toVersion: "9.8.7", toInclusive: false},
			want:    want{fromString: "1.0.0", fromInclusive: false, toString: "9.8.7", toInclusive: false},
			wantErr: false,
		},
		{
			name:    "to_major",
			args:    args{fromVersion: "1.2.3", fromInclusive: false, toVersion: "9", toInclusive: false},
			want:    want{fromString: "1.2.3", fromInclusive: false, toString: "9.0.0", toInclusive: false},
			wantErr: false,
		},
		{
			name:    "both_major",
			args:    args{fromVersion: "1", fromInclusive: false, toVersion: "9", toInclusive: false},
			want:    want{fromString: "1.0.0", fromInclusive: false, toString: "9.0.0", toInclusive: false},
			wantErr: false,
		},
		{
			name:    "from_invalid_semantic_version",
			args:    args{fromVersion: "1.2.3.4", fromInclusive: false, toVersion: "9.8.7", toInclusive: false},
			want:    want{},
			wantErr: true,
		},
		{
			name:    "to_invalid_semantic_version",
			args:    args{fromVersion: "1.2.3", fromInclusive: false, toVersion: "9.8.7.6", toInclusive: false},
			want:    want{},
			wantErr: true,
		},
		{
			name:    "both_invalid_semantic_version",
			args:    args{fromVersion: "1.2.3.4", fromInclusive: false, toVersion: "9.8.7.6", toInclusive: false},
			want:    want{},
			wantErr: true,
		},
		{
			name:    "from_invalid_version",
			args:    args{fromVersion: "informational", fromInclusive: false, toVersion: "9.8.7", toInclusive: false},
			want:    want{},
			wantErr: true,
		},
		{
			name:    "to_invalid_version",
			args:    args{fromVersion: "1.2.3", fromInclusive: false, toVersion: "informational", toInclusive: false},
			want:    want{},
			wantErr: true,
		},
		{
			name:    "both_invalid_version",
			args:    args{fromVersion: "informational", fromInclusive: false, toVersion: "informational", toInclusive: false},
			want:    want{},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewRange(
				tc.args.fromVersion,
				tc.args.fromInclusive,
				tc.args.toVersion,
				tc.args.toInclusive,
			)

			if tc.wantErr {
				if err == nil {
					t.Errorf("NewRange() error = %v, wantErr %v", err, tc.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewRange() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if got.fromInclusive != tc.want.fromInclusive {
				t.Errorf(
					"NewRange().fromInclusive got = %v, want %v",
					got.fromInclusive,
					tc.want.fromInclusive,
				)
			}

			if got.fromVersion.String() != tc.want.fromString {
				t.Errorf(
					"NewRange().fromVersion.String() got = %v, want %v",
					got.fromVersion.String(),
					tc.want.fromString,
				)
			}

			if got.toInclusive != tc.want.toInclusive {
				t.Errorf(
					"NewRange().toInclusive got = %v, want %v",
					got.toInclusive,
					tc.want.toInclusive,
				)
			}

			if got.toVersion.String() != tc.want.toString {
				t.Errorf(
					"NewRange().toVersion.String() got = %v, want %v",
					got.toVersion.String(),
					tc.want.toString,
				)
			}

			if got.fromWildcard {
				t.Errorf(
					"NewRange().fromWildcard got = %v, want %v",
					got.fromWildcard,
					false,
				)
			}

			if got.toWildcard {
				t.Errorf(
					"NewRange().toWildcard got = %v, want %v",
					got.toWildcard,
					false,
				)
			}
		})
	}
}

func TestRange_String(t *testing.T) {
	type fields struct {
		fromVersion   string
		fromInclusive bool
		toVersion     string
		toInclusive   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "from_inclusive",
			fields: fields{fromVersion: "1.2.3", fromInclusive: true, toVersion: "9.8.7", toInclusive: false},
			want:   ">=1.2.3,<9.8.7",
		},
		{
			name:   "to_inclusive",
			fields: fields{fromVersion: "1.2.3", fromInclusive: false, toVersion: "9.8.7", toInclusive: true},
			want:   ">1.2.3,<=9.8.7",
		},
		{
			name:   "both_inclusive",
			fields: fields{fromVersion: "1.2.3", fromInclusive: true, toVersion: "9.8.7", toInclusive: true},
			want:   ">=1.2.3,<=9.8.7",
		},
		{
			name:   "both_non-inclusive",
			fields: fields{fromVersion: "1.2.3", fromInclusive: false, toVersion: "9.8.7", toInclusive: false},
			want:   ">1.2.3,<9.8.7",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := NewRange(
				tc.fields.fromVersion,
				tc.fields.fromInclusive,
				tc.fields.toVersion,
				tc.fields.toInclusive,
			)

			got := r.String()

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("String() got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRange_String_upper_bounded(t *testing.T) {
	type fields struct {
		version   string
		inclusive bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "inclusive",
			fields: fields{version: "1.2.3", inclusive: true},
			want:   "<=1.2.3",
		},
		{
			name:   "non-inclusive",
			fields: fields{version: "1.2.3", inclusive: false},
			want:   "<1.2.3",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := NewUpperBoundedRange(tc.fields.version, tc.fields.inclusive)

			got := r.String()

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("String() got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRange_String_lower_bounded(t *testing.T) {
	type fields struct {
		version   string
		inclusive bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "inclusive",
			fields: fields{version: "1.2.3", inclusive: true},
			want:   ">=1.2.3",
		},
		{
			name:   "non-inclusive",
			fields: fields{version: "1.2.3", inclusive: false},
			want:   ">1.2.3",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := NewLowerBoundedRange(tc.fields.version, tc.fields.inclusive)

			got := r.String()

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("String() got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRange_String_wildcard(t *testing.T) {
	r, _ := NewWildcardRange()

	got := r.String()

	if got != "*" {
		t.Errorf("String() got = %v, want %v", got, "*")
	}
}
