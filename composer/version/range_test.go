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

package version

import (
	"reflect"
	"testing"
)

var (
	v1 = &Version{1, 0, 0, 0}
	v2 = &Version{2, 0, 0, 0}
	v3 = &Version{3, 0, 0, 0}
	v4 = &Version{4, 0, 0, 0}
	v5 = &Version{5, 0, 0, 0}

	unbounded = Range{}

	r1i2  = Range{v1, true, v2, false}
	r1i2i = Range{v1, true, v2, true}
	r1i3  = Range{v1, true, v3, false}
	r1i3i = Range{v1, true, v3, true}
	r2i3  = Range{v2, true, v3, false}
	r2i3i = Range{v2, true, v3, true}
	r12   = Range{v1, false, v2, false}
	r12i  = Range{v1, false, v2, true}
	r13   = Range{v1, false, v3, false}
	r13i  = Range{v1, false, v3, true}
	r14   = Range{v1, false, v4, false}
	r15   = Range{v1, false, v5, false}
	r23   = Range{v2, false, v3, false}
	r24   = Range{v2, false, v4, false}
	r34   = Range{v3, false, v4, false}
	r35   = Range{v3, false, v5, false}

	l1 = Range{nil, false, v1, false}
	l2 = Range{nil, false, v2, false}
	l3 = Range{nil, false, v3, false}
	l4 = Range{nil, false, v4, false}

	g1 = Range{v1, false, nil, false}
	g2 = Range{v2, false, nil, false}
	g3 = Range{v3, false, nil, false}
	g4 = Range{v4, false, nil, false}
)

func TestNewRange(t *testing.T) {
	tests := []struct {
		optFns []RangeOptFn
		want   Range
	}{
		{
			[]RangeOptFn{},
			Range{nil, false, nil, false},
		},
		{
			[]RangeOptFn{
				WithInclusiveCeiling(v1),
			},
			Range{nil, false, v1, true},
		},
		{
			[]RangeOptFn{
				WithNonInclusiveCeiling(v1),
			},
			Range{nil, false, v1, false},
		},
		{
			[]RangeOptFn{
				WithoutCeiling(),
			},
			Range{nil, false, nil, false},
		},
		{
			[]RangeOptFn{
				WithInclusiveCeiling(v1),
				WithoutCeiling(),
			},
			Range{nil, false, nil, false},
		},
		{
			[]RangeOptFn{
				WithNonInclusiveCeiling(v1),
				WithoutCeiling(),
			},
			Range{nil, false, nil, false},
		},
		{
			[]RangeOptFn{
				WithInclusiveFloor(v1),
			},
			Range{v1, true, nil, false},
		},
		{
			[]RangeOptFn{
				WithNonInclusiveFloor(v1),
			},
			Range{v1, false, nil, false},
		},
		{
			[]RangeOptFn{
				WithoutFloor(),
			},
			Range{nil, false, nil, false},
		},
		{
			[]RangeOptFn{
				WithInclusiveFloor(v1),
				WithoutFloor(),
			},
			Range{nil, false, nil, false},
		},
		{
			[]RangeOptFn{
				WithNonInclusiveFloor(v1),
				WithoutFloor(),
			},
			Range{nil, false, nil, false},
		},
		{
			[]RangeOptFn{
				WithInclusiveCeiling(v1),
				WithInclusiveFloor(v2),
			},
			Range{v2, true, v1, true},
		},
		{
			[]RangeOptFn{
				WithInclusiveCeiling(v1),
				WithNonInclusiveFloor(v2),
			},
			Range{v2, false, v1, true},
		},
		{
			[]RangeOptFn{
				WithInclusiveCeiling(v1),
				WithoutFloor(),
			},
			Range{nil, false, v1, true},
		},
		{
			[]RangeOptFn{
				WithNonInclusiveCeiling(v1),
				WithInclusiveFloor(v2),
			},
			Range{v2, true, v1, false},
		},
		{
			[]RangeOptFn{
				WithNonInclusiveCeiling(v1),
				WithNonInclusiveFloor(v2),
			},
			Range{v2, false, v1, false},
		},
		{
			[]RangeOptFn{
				WithNonInclusiveCeiling(v1),
				WithoutFloor(),
			},
			Range{nil, false, v1, false},
		},
		{
			[]RangeOptFn{
				WithoutCeiling(),
				WithInclusiveFloor(v1),
			},
			Range{v1, true, nil, false},
		},
		{
			[]RangeOptFn{
				WithoutCeiling(),
				WithNonInclusiveFloor(v1),
			},
			Range{v1, false, nil, false},
		},
		{
			[]RangeOptFn{
				WithoutCeiling(),
				WithoutFloor(),
			},
			Range{nil, false, nil, false},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if reflect.DeepEqual(v1, v2) {
				t.Errorf("bad implmentation v1 == v2")
			}

			if got := NewRange(tt.optFns...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRange_String(t *testing.T) {
	v1234 := &Version{1, 2, 3, 4}
	v9876 := &Version{9, 8, 7, 6}

	type fields struct {
		from          *Version
		fromInclusive bool
		to            *Version
		toInclusive   bool
	}
	tests := []struct {
		fields fields
		want   string
	}{
		{fields{}, "*"},
		{fields{v1234, false, nil, false}, ">1.2.3.4"},
		{fields{v1234, true, nil, false}, ">=1.2.3.4"},
		{fields{v1234, false, nil, true}, ">1.2.3.4"},
		{fields{v1234, true, nil, true}, ">=1.2.3.4"},
		{fields{nil, false, v1234, false}, "<1.2.3.4"},
		{fields{nil, false, v1234, true}, "<=1.2.3.4"},
		{fields{nil, true, v1234, false}, "<1.2.3.4"},
		{fields{nil, true, v1234, true}, "<=1.2.3.4"},
		{fields{v1234, false, v9876, false}, ">1.2.3.4 <9.8.7.6"},
		{fields{v1234, false, v9876, true}, ">1.2.3.4 <=9.8.7.6"},
		{fields{v1234, true, v9876, false}, ">=1.2.3.4 <9.8.7.6"},
		{fields{v1234, true, v9876, true}, ">=1.2.3.4 <=9.8.7.6"},
		{fields{v1234, true, v1234, true}, "1.2.3.4"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			r := Range{
				from:          tt.fields.from,
				fromInclusive: tt.fields.fromInclusive,
				to:            tt.fields.to,
				toInclusive:   tt.fields.toInclusive,
			}
			if got := r.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstraint_String(t *testing.T) {
	tests := []struct {
		name string
		c    Constraint
		want string
	}{
		{
			"empty",
			[]Range{},
			"",
		},
		{
			"single",
			[]Range{r12},
			">1 <2",
		},
		{
			"no overlap",
			[]Range{r12, r23, r34},
			">1 <2||>2 <3||>3 <4",
		},
		{
			"without floor",
			[]Range{l1, l2},
			"<2",
		},
		{
			"without floor",
			[]Range{l2, l3, l1},
			"<3",
		},
		{
			"without floor",
			[]Range{l4, l3, l2, l1},
			"<4",
		},
		{
			"without celling",
			[]Range{g1, g2},
			">1",
		},
		{
			"without celling",
			[]Range{g2, g3, g1},
			">1",
		},
		{
			"without celling",
			[]Range{g4, g3, g2, g1},
			">1",
		},
		{
			"overlap",
			[]Range{r13, r24},
			">1 <4",
		},
		{
			"overlap",
			[]Range{r24, r13, r35},
			">1 <5",
		},
		{
			"touching",
			[]Range{r12i, r23},
			">1 <3",
		},
		{
			"touching",
			[]Range{r12, r2i3},
			">1 <3",
		},
		{
			"touching",
			[]Range{r1i2i, r23},
			">=1 <3",
		},
		{
			"touching",
			[]Range{r12, r2i3i},
			">1 <=3",
		},
		{
			"touching",
			[]Range{r1i2, r2i3i},
			">=1 <=3",
		},
		{
			"included",
			[]Range{r14, r23},
			">1 <4",
		},
		{
			"included",
			[]Range{r23, r14},
			">1 <4",
		},
		{
			"unbounded",
			[]Range{unbounded},
			"*",
		},
		{
			"unbounded",
			[]Range{unbounded, unbounded},
			"*",
		},
		{
			"unbounded",
			[]Range{r12, unbounded},
			"*",
		},
		{
			"unbounded",
			[]Range{unbounded, r12},
			"*",
		},
		{
			"unbounded",
			[]Range{r14, r23, unbounded, r13, r24},
			"*",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_or(t *testing.T) {
	tests := []struct {
		name string
		args []Range
		want []Range
	}{
		{
			"empty",
			[]Range{},
			[]Range{},
		},
		{
			"single",
			[]Range{r12},
			[]Range{r12},
		},
		{
			"no overlap",
			[]Range{r12, r23, r34},
			[]Range{r12, r23, r34},
		},
		{
			"without floor",
			[]Range{l1, l2},
			[]Range{l2},
		},
		{
			"without floor",
			[]Range{l2, l3, l1},
			[]Range{l3},
		},
		{
			"without floor",
			[]Range{l4, l3, l2, l1},
			[]Range{l4},
		},
		{
			"without celling",
			[]Range{g1, g2},
			[]Range{g1},
		},
		{
			"without celling",
			[]Range{g2, g3, g1},
			[]Range{g1},
		},
		{
			"without celling",
			[]Range{g4, g3, g2, g1},
			[]Range{g1},
		},
		{
			"overlap",
			[]Range{r13, r24},
			[]Range{r14},
		},
		{
			"overlap",
			[]Range{r24, r13, r35},
			[]Range{r15},
		},
		{
			"touching",
			[]Range{r12i, r23},
			[]Range{r13},
		},
		{
			"touching",
			[]Range{r12, r2i3},
			[]Range{r13},
		},
		{
			"touching",
			[]Range{r1i2i, r23},
			[]Range{r1i3},
		},
		{
			"touching",
			[]Range{r12, r2i3i},
			[]Range{r13i},
		},
		{
			"touching",
			[]Range{r1i2, r2i3i},
			[]Range{r1i3i},
		},
		{
			"included",
			[]Range{r14, r23},
			[]Range{r14},
		},
		{
			"included",
			[]Range{r23, r14},
			[]Range{r14},
		},
		{
			"unbounded",
			[]Range{unbounded},
			[]Range{unbounded},
		},
		{
			"unbounded",
			[]Range{unbounded, unbounded},
			[]Range{unbounded},
		},
		{
			"unbounded",
			[]Range{r12, unbounded},
			[]Range{unbounded},
		},
		{
			"unbounded",
			[]Range{unbounded, r12},
			[]Range{unbounded},
		},
		{
			"unbounded",
			[]Range{r14, r23, unbounded, r13, r24},
			[]Range{unbounded},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := or(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("or() = %v, want %v", got, tt.want)
			}
		})
	}
}