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
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	regex = `^v?[[:blank:]]?(?P<major>\d+)(\.(?P<minor>\d+)(\.(?P<patch>\d+)(\.(?P<revision>\d+)?)?)?)?`
)

var (
	ErrMalformedVersionString = errors.New("malformed version string")
)

type Version struct {
	major    uint64
	minor    uint64
	patch    uint64
	revision uint64
}

// NewVersion parses a given string and returns a pointer of [Version] or
// an error if unable to parse the version. If the version is not composer
// compilable it attempts to convert it.
func NewVersion(v string) (*Version, error) {
	r := regexp.MustCompile(regex)

	if !r.MatchString(v) {
		return nil, fmt.Errorf("unable to parse %q: %w", v, ErrMalformedVersionString)
	}

	ver := &Version{}

	m := r.FindStringSubmatch(v)
	for i, name := range r.SubexpNames() {
		n, err := strconv.ParseUint(m[i], 10, 0)
		if err != nil {
			continue
		}

		switch name {
		case "major":
			ver.major = n
		case "minor":
			ver.minor = n
		case "patch":
			ver.patch = n
		case "revision":
			ver.revision = n
		}
	}

	return ver, nil
}

// normalize returns the normalized formatting of the version. It fills in any
// missing .MINOR or .PATCH or .REVISION. Two versions are compared equal only
// if their normalized formatting are identical strings.
func (v Version) normalize() string {
	return fmt.Sprintf("%d.%d.%d.%d", v.major, v.minor, v.patch, v.revision)
}

func (v Version) String() string {
	s := v.normalize()

	s, _ = strings.CutSuffix(s, ".0") // revision
	s, _ = strings.CutSuffix(s, ".0") // patch
	s, _ = strings.CutSuffix(s, ".0") // minor

	return s
}

// compare this version to another one. It returns -1, 0 or 1 if the version
// is smaller, equal or larger than the other version respectively.
func (v Version) compare(o Version) int {
	if d := compareUint(v.major, o.major); d != 0 {
		return d
	}
	if d := compareUint(v.minor, o.minor); d != 0 {
		return d
	}
	if d := compareUint(v.patch, o.patch); d != 0 {
		return d
	}

	return compareUint(v.revision, o.revision)
}

func compareUint(i, j uint64) int {
	switch {
	case i < j:
		return -1
	case i > j:
		return 1
	default:
		return 0
	}
}

func (v Version) equalTo(o Version) bool {
	return v.compare(o) == 0
}

func (v Version) greaterThan(o Version) bool {
	return v.compare(o) > 0
}

func (v Version) lessThan(o Version) bool {
	return v.compare(o) < 0
}