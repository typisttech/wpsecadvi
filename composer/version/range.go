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
	"fmt"
	"slices"
	"strings"
)

type Range struct {
	from          *Version // TODO: Can we use value instead?
	fromInclusive bool
	to            *Version // TODO: Can we use value instead?
	toInclusive   bool
}

type RangeOptFn func(*Range)

func WithInclusiveCeiling(v *Version) RangeOptFn {
	return func(r *Range) {
		r.to = v
		r.toInclusive = true
	}
}

func WithNonInclusiveCeiling(v *Version) RangeOptFn {
	return func(r *Range) {
		r.to = v
		r.toInclusive = false
	}
}

func WithoutCeiling() RangeOptFn {
	return func(r *Range) {
		r.to = nil
		r.toInclusive = false
	}
}

func WithInclusiveFloor(v *Version) RangeOptFn {
	return func(r *Range) {
		r.from = v
		r.fromInclusive = true
	}
}

func WithNonInclusiveFloor(v *Version) RangeOptFn {
	return func(r *Range) {
		r.from = v
		r.fromInclusive = false
	}
}

func WithoutFloor() RangeOptFn {
	return func(r *Range) {
		r.from = nil
		r.fromInclusive = false
	}
}

func NewRange(optFns ...RangeOptFn) Range {
	r := &Range{}

	for _, optFn := range optFns {
		optFn(r)
	}

	// TODO: Validate r.from <= r.to
	// TODO: Validate when r.from == r.to, then r.fromInclusive == true && r.toInclusive == true

	return *r
}

func (r Range) String() string {
	if r.from == nil && r.to == nil {
		return "*"
	}

	if r.fromInclusive && r.toInclusive && r.from != nil && r.from == r.to {
		return r.from.String()
	}

	from := ""
	if r.from != nil {
		op := ">"
		if r.fromInclusive {
			op = ">="
		}

		from = fmt.Sprintf("%s%s", op, r.from)
	}

	to := ""
	if r.to != nil {
		op := "<"
		if r.toInclusive {
			op = "<="
		}

		to = fmt.Sprintf("%s%s", op, r.to)
	}

	str := strings.Join([]string{from, to}, " ")
	return strings.Trim(str, " ")
}

// Constraint represents a slice of [Range] grouped together with logical OR.
type Constraint []Range

func (c Constraint) String() string {
	rs := or(c...)

	ss := make([]string, 0, len(rs))
	for _, r := range rs {
		ss = append(ss, r.String())
	}

	// For easier testing assertions.
	slices.Sort(ss)

	return strings.Join(ss, "||")
}

func or(rs ...Range) []Range {
	if len(rs) == 0 {
		return []Range{}
	}

	f, rs := rs[0], rs[1:]
	result := []Range{f}

	for i, r := range rs {
		for j, s := range result {
			if m, ok := orTwo(r, s); ok {
				result[j] = m
				result = append(result, rs[i+1:]...)
				return or(result...)
			}
		}
		result = append(result, r)
	}

	return result
}

func orTwo(a Range, b Range) (Range, bool) {
	if !overlap(a, b) {
		return Range{}, false
	}

	if a.from == nil && a.to == nil {
		return a, true
	}
	if b.from == nil && b.to == nil {
		return b, true
	}
	if a.String() == b.String() {
		return a, true
	}

	// Both without celling, take the lesser from
	//	    |<-a->
	//	|<---b--->
	if a.to == nil && b.to == nil {
		from := a.from
		fromInclusive := a.fromInclusive

		if b.from.lessThan(*a.from) {
			from = b.from
			fromInclusive = b.fromInclusive
		}
		if a.from.equalTo(*b.from) {
			fromInclusive = a.fromInclusive && b.fromInclusive
		}

		return Range{from, fromInclusive, nil, false}, true
	}

	// Both without floor
	//	<-a->|
	//	<---b--->|
	if a.from == nil && b.from == nil {
		to := a.to
		toInclusive := a.toInclusive

		if b.to.greaterThan(*a.to) {
			to = b.to
			toInclusive = b.toInclusive
		}
		if a.to.equalTo(*b.to) {
			toInclusive = a.toInclusive && b.toInclusive
		}

		return Range{nil, false, to, toInclusive}, true
	}

	// Ensure a has a lower from
	//	  |<-a->|
	//	      |<---b--->|
	// Or,
	//	  |<-a->|
	//	        |<---b--->|
	// Or,
	//	  |<-a->|
	//	           |<--- b --->|
	// Or,
	//	  |<--------a-------->|
	//	        |<---b--->|
	if a.from.greaterThan(*b.from) {
		a, b = b, a
	}

	from := a.from
	fromInclusive := a.fromInclusive
	if a.from.equalTo(*b.from) {
		fromInclusive = a.fromInclusive || b.fromInclusive
	}

	to := b.to
	toInclusive := b.toInclusive
	if a.to.greaterThan(*b.to) {
		to = a.to
	}
	if a.to.equalTo(*b.to) {
		toInclusive = a.toInclusive || b.toInclusive
	}

	return Range{from, fromInclusive, to, toInclusive}, true
}

func overlap(a Range, b Range) bool {
	if a.from == nil && a.to == nil {
		return true
	}
	if b.from == nil && b.to == nil {
		return true
	}
	if a.String() == b.String() {
		return true
	}

	// Both without celling
	//	    |<-a->
	//	|<---b--->
	if a.to == nil && b.to == nil {
		return true
	}

	// Both without floor
	//	<-a->|
	//	<---b--->|
	if a.from == nil && b.from == nil {
		return true
	}

	// Ensure a has a lower from
	//	  |<-a->|
	//	      |<---b--->|
	// Or,
	//	  |<-a->|
	//	        |<---b--->|
	// Or,
	//	  |<-a->|
	//	           |<--- b --->|
	// Or,
	//	  |<--------a-------->|
	//	        |<---b--->|
	if a.from.greaterThan(*b.from) {
		a, b = b, a
	}

	return a.to == nil ||
		a.to.greaterThan(*b.from) ||
		(a.to.equalTo(*b.from) && (a.toInclusive || b.fromInclusive))
}