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
	"fmt"
	"github.com/Masterminds/semver/v3"
	"golang.org/x/exp/slices"
	"strings"
)

type Constraints struct {
	ranges []*Range
}

func (cs *Constraints) Or(r *Range) {
	rs := make([]*Range, 0, len(cs.ranges)+1)
	for _, csr := range cs.ranges {
		if csr.cover(r) {
			return
		}

		if r.cover(csr) {
			continue
		}

		rs = append(rs, csr)
	}

	cs.ranges = append(rs, r)
}

func (cs Constraints) String() string {
	rs := make([]string, 0, len(cs.ranges))
	for _, r := range cs.ranges {
		rStr := r.String()
		if rStr == "" {
			continue
		}
		rs = append(rs, r.String())
	}

	// For easier testing assertions.
	slices.Sort(rs)

	return strings.Join(rs, "|")
}

type Range struct {
	fromVersion   *semver.Version
	fromInclusive bool
	fromWildcard  bool
	toVersion     *semver.Version
	toInclusive   bool
	toWildcard    bool
}

func NewWildcardRange() (*Range, error) {
	return &Range{
		fromWildcard: true,
		toWildcard:   true,
	}, nil
}

func NewUpperBoundedRange(version string, inclusive bool) (*Range, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}

	return &Range{
		fromWildcard: true,
		toVersion:    v,
		toInclusive:  inclusive,
	}, nil
}

func NewLowerBoundedRange(version string, inclusive bool) (*Range, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}

	return &Range{
		fromVersion:   v,
		fromInclusive: inclusive,
		toWildcard:    true,
	}, nil
}

func NewRange(fromVersion string, fromInclusive bool, toVersion string, toInclusive bool) (*Range, error) {
	fromV, err := semver.NewVersion(fromVersion)
	if err != nil {
		return nil, err
	}

	toV, err := semver.NewVersion(toVersion)
	if err != nil {
		return nil, err
	}

	return &Range{
		fromVersion:   fromV,
		fromInclusive: fromInclusive,
		toVersion:     toV,
		toInclusive:   toInclusive,
	}, nil
}

func (r Range) String() string {
	if r.fromWildcard && r.toWildcard {
		return "*"
	}

	from := ""
	if !r.fromWildcard {
		op := ">"
		if r.fromInclusive {
			op = ">="
		}

		from = fmt.Sprintf("%s%s", op, r.fromVersion)
	}

	to := ""
	if !r.toWildcard {
		op := "<"
		if r.toInclusive {
			op = "<="
		}

		to = fmt.Sprintf("%s%s", op, r.toVersion)
	}

	csStr := strings.Join([]string{from, to}, ",")
	return strings.Trim(csStr, ",")
}

func (r *Range) cover(other *Range) bool {
	fromCovered := r.fromWildcard ||
		(r.fromVersion != nil && other.fromVersion != nil && r.fromVersion.LessThan(other.fromVersion)) ||
		(r.fromVersion != nil && other.fromVersion != nil && (r.fromInclusive || !other.fromInclusive) && r.fromVersion.Equal(other.fromVersion))

	toCovered := r.toWildcard ||
		(r.toVersion != nil && other.toVersion != nil && r.toVersion.GreaterThan(other.toVersion)) ||
		(r.toVersion != nil && other.toVersion != nil && (r.toInclusive || !other.toInclusive) && r.toVersion.Equal(other.toVersion))

	return fromCovered && toCovered
}
