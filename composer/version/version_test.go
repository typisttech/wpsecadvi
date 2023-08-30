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
	"testing"
)

// tests is modified from `composer/semver`'s test.
// https://github.com/composer/semver/blob/fa1ec24f0ab1efe642671ec15c51a3ab879f59bf/tests/VersionParserTest.php#L65-L134
// https://github.com/composer/semver/blob/fa1ec24f0ab1efe642671ec15c51a3ab879f59bf/tests/VersionParserTest.php#L150-L174
var tests = []struct {
	name string
	in   string
	out  string
}{
	// From `VersionParserTest::successfulNormalizedVersions()`
	{"none", "1.0.0", "1.0.0.0"},
	{"none/2", "1.2.3.4", "1.2.3.4"},
	//{"parses state", "1.0.0RC1dev", "1.0.0.0-RC1-dev"},
	//{"CI parsing", "1.0.0-rC15-dev", "1.0.0.0-RC15-dev"},
	//{"delimiters", "1.0.0.RC.15-dev", "1.0.0.0-RC15-dev"},
	//{"RC uppercase", "1.0.0-rc1", "1.0.0.0-RC1"},
	//{"patch replace", "1.0.0.pl3-dev", "1.0.0.0-patch3-dev"},
	//{"forces w.x.y.z", "1.0-dev", "1.0.0.0-dev"},
	{"forces w.x.y.z/2", "0", "0.0.0.0"},
	//{"parses long", "10.4.13-beta", "10.4.13.0-beta"},
	//{"parses long/2", "10.4.13beta2", "10.4.13.0-beta2"},
	//{"parses long/semver", "10.4.13beta.2", "10.4.13.0-beta2"},
	//{"parses long/semver2", "v1.13.11-beta.0", "1.13.11.0-beta0"},
	//{"parses long/semver3", "1.13.11.0-beta0", "1.13.11.0-beta0"},
	//{"expand shorthand", "10.4.13-b", "10.4.13.0-beta"},
	//{"expand shorthand/2", "10.4.13-b5", "10.4.13.0-beta5"},
	{"strips leading v", "v1.0.0", "1.0.0.0"},
	//{"parses dates y-m as classical", "2010.01", "2010.01.0.0"},
	//{"parses dates w/ . as classical", "2010.01.02", "2010.01.02.0"},
	{"parses dates y.m.Y as classical", "2010.1.555", "2010.1.555.0"},
	{"parses dates y.m.Y/2 as classical", "2010.10.200", "2010.10.200.0"},
	//{"strips v/datetime", "v20100102", "20100102"},
	//{"parses dates w/ -", "2010-01-02", "2010.01.02"},
	//{"parses dates w/ .", "2012.06.07", "2012.06.07.0"},
	//{"parses numbers", "2010-01-02.5", "2010.01.02.5"},
	{"parses dates y.m.Y", "2010.1.555", "2010.1.555.0"},
	//{"parses datetime", "20100102-203040", "20100102.203040"},
	//{"parses date dev", "20100102.x-dev", "20100102.9999999.9999999.9999999-dev"},
	//{"parses datetime dev", "20100102.203040.x-dev", "20100102.203040.9999999.9999999-dev"},
	//{"parses dt+number", "20100102203040-10", "20100102203040.10"},
	//{"parses dt+patch", "20100102-203040-p1", "20100102.203040-patch1"},
	//{"parses dt Ym", "201903.0", "201903.0"},
	//{"parses dt Ym dev", "201903.x-dev", "201903.9999999.9999999.9999999-dev"},
	//{"parses dt Ym+patch", "201903.0-p2", "201903.0-patch2"},
	//{"parses master", "dev-master", "dev-master"},
	//{"parses master w/o dev", "master", "dev-master"},
	//{"parses trunk", "dev-trunk", "dev-trunk"},
	//{"parses branches", "1.x-dev", "1.9999999.9999999.9999999-dev"},
	//{"parses arbitrary", "dev-feature-foo", "dev-feature-foo"},
	//{"parses arbitrary/2", "DEV-FOOBAR", "dev-FOOBAR"},
	//{"parses arbitrary/3", "dev-feature/foo", "dev-feature/foo"},
	//{"parses arbitrary/4", "dev-feature+issue-1", "dev-feature+issue-1"},
	//{"ignores aliases", "dev-master as 1.0.0", "dev-master"},
	//{"ignores aliases/2", "dev-load-varnish-only-when-used as ^2.0", "dev-load-varnish-only-when-used"},
	//{"ignores aliases/3", "dev-load-varnish-only-when-used@dev as ^2.0@dev", "dev-load-varnish-only-when-used"},
	{"ignores stability", "1.0.0+foo@dev", "1.0.0.0"},
	//{"ignores stability/2", "dev-load-varnish-only-when-used@stable", "dev-load-varnish-only-when-used"},
	//{"semver metadata/2", "1.0.0-beta.5+foo", "1.0.0.0-beta5"},
	{"semver metadata/3", "1.0.0+foo", "1.0.0.0"},
	//{"semver metadata/4", "1.0.0-alpha.3.1+foo", "1.0.0.0-alpha3.1"},
	//{"semver metadata/5", "1.0.0-alpha2.1+foo", "1.0.0.0-alpha2.1"},
	//{"semver metadata/6", "1.0.0-alpha-2.1-3+foo", "1.0.0.0-alpha2.1-3"},
	// not supported for BC "semver metadata/7", "1.0.0-0.3.7", "1.0.0.0-0.3.7"},
	// not supported for BC "semver metadata/8", "1.0.0-x.7.z.92", "1.0.0.0-x.7.z.92"},
	{"metadata w/ alias", "1.0.0+foo as 2.0", "1.0.0.0"},
	//{"keep zero-padding", "00.01.03.04", "00.01.03.04"},
	//{"keep zero-padding/2", "000.001.003.004", "000.001.003.004"},
	//{"keep zero-padding/3", "0.000.103.204", "0.000.103.204"},
	//{"keep zero-padding/4", "0700", "0700.0.0.0"},
	//{"keep zero-padding/5", "041.x-dev", "041.9999999.9999999.9999999-dev"},
	//{"keep zero-padding/6", "dev-041.003", "dev-041.003"},
	//{"dev with mad name", "dev-1.0.0-dev<1.0.5-dev", "dev-1.0.0-dev<1.0.5-dev"},
	//{"dev prefix with spaces", "dev-foo bar", "dev-foo bar"},
	{"space padding", " 1.0.0", "1.0.0.0"},
	{"space padding/2", "1.0.0 ", "1.0.0.0"},

	// From `VersionParserTest::failingNormalizedVersions()`
	{"empty ", "", ""},
	{"invalid chars", "a", ""},
	//{"invalid type", "1.0.0-meh", ""},
	//{"too many bits", "1.0.0.0.0", ""},
	{"non-dev arbitrary", "feature-foo", ""},
	//{"metadata w/ space", "1.0.0+foo bar", ""},
	//{"maven style release", "1.0.1-SNAPSHOT", ""},
	//{"dev with less than", "1.0.0<1.0.5-dev", ""},
	//{"dev with less than/2", "1.0.0-dev<1.0.5-dev", ""},
	{"dev suffix with spaces", "foo bar-dev", ""},
	//{"any with spaces", "1.0 .2", ""},
	{"no version, no alias", " as ", ""},
	{"no version, only alias", " as 1.2", ""},
	{"just an operator", "^", ""},
	{"just an operator/2", "^8 || ^", ""},
	{"just an operator/3", "~", ""},
	{"just an operator/4", "~1 ~", ""},
	{"constraint", "~1", ""},
	{"constraint/2", "^1", ""},
	//{"constraint/3", "1.*", ""},
}

func TestNewVersion(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.in)

			if tt.out == "" {
				if err == nil {
					t.Errorf("New() error = nil, want error")
				}

				if !errors.Is(err, ErrMalformedVersionString) {
					t.Errorf("New() error = %v, want ErrMalformedVersionString", err)
				}

				return
			}

			if err != nil {
				t.Errorf("New() error = %v, want nil", err)
				return
			}
			if got.normalize() != tt.out {
				t.Errorf("New() got = %v, want %v", got, tt.out)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	type fields struct {
		major    uint64
		minor    uint64
		patch    uint64
		revision uint64
	}
	tests := []struct {
		fields fields
		want   string
	}{
		// No zero.
		{fields{1, 2, 3, 4}, "1.2.3.4"},
		//// 1 zero.
		{fields{1, 2, 3, 0}, "1.2.3"},
		{fields{1, 2, 0, 4}, "1.2.0.4"},
		{fields{1, 0, 3, 4}, "1.0.3.4"},
		{fields{0, 2, 3, 4}, "0.2.3.4"},
		//// 2 zeros.
		{fields{1, 2, 0, 0}, "1.2"},
		{fields{1, 0, 3, 0}, "1.0.3"},
		{fields{0, 2, 3, 0}, "0.2.3"},
		{fields{1, 0, 0, 4}, "1.0.0.4"},
		{fields{0, 2, 0, 4}, "0.2.0.4"},
		{fields{0, 0, 3, 4}, "0.0.3.4"},
		//// 3 zeros.
		{fields{1, 0, 0, 0}, "1"},
		{fields{0, 2, 0, 0}, "0.2"},
		{fields{0, 0, 3, 0}, "0.0.3"},
		{fields{0, 0, 0, 4}, "0.0.0.4"},
		// 4 zeros.
		{fields{0, 0, 0, 0}, "0"},
		//
		{fields{10, 2, 3, 4}, "10.2.3.4"},
		{fields{1, 20, 3, 4}, "1.20.3.4"},
		{fields{1, 2, 30, 4}, "1.2.30.4"},
		{fields{1, 2, 3, 40}, "1.2.3.40"},
		{fields{10, 20, 30, 40}, "10.20.30.40"},
	}
	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%d.%d.%d.%d", tt.fields.major, tt.fields.minor, tt.fields.patch, tt.fields.revision),
			func(t *testing.T) {
				v := Version{
					major:    tt.fields.major,
					minor:    tt.fields.minor,
					patch:    tt.fields.patch,
					revision: tt.fields.revision,
				}
				if got := v.String(); got != tt.want {
					t.Errorf("String() = %v, want %v", got, tt.want)
				}
			})
	}
}

func TestVersion_compare(t *testing.T) {
	tests := []struct {
		v    Version
		o    Version
		want int
	}{
		{Version{1, 2, 3, 4}, Version{1, 2, 3, 4}, 0},
		{Version{2, 2, 3, 4}, Version{1, 2, 3, 4}, 1},
		{Version{1, 3, 3, 4}, Version{1, 2, 3, 4}, 1},
		{Version{1, 2, 4, 4}, Version{1, 2, 3, 4}, 1},
		{Version{1, 2, 3, 5}, Version{1, 2, 3, 4}, 1},
		{Version{0, 2, 3, 4}, Version{1, 2, 3, 4}, -1},
		{Version{1, 1, 3, 4}, Version{1, 2, 3, 4}, -1},
		{Version{1, 2, 2, 4}, Version{1, 2, 3, 4}, -1},
		{Version{1, 2, 3, 3}, Version{1, 2, 3, 4}, -1},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.v.compare(tt.o); got != tt.want {
				t.Errorf("'%s'.compare('%s') = %v, want %v", tt.v, tt.o, got, tt.want)
			}
		})
	}
}

func TestVersion_equalTo(t *testing.T) {
	tests := []struct {
		v    Version
		o    Version
		want bool
	}{
		{Version{0, 0, 0, 0}, Version{0, 0, 0, 0}, true},
		{Version{1, 2, 3, 4}, Version{1, 2, 3, 4}, true},
		{Version{2, 2, 3, 4}, Version{1, 2, 3, 4}, false},
		{Version{1, 3, 3, 4}, Version{1, 2, 3, 4}, false},
		{Version{1, 2, 4, 4}, Version{1, 2, 3, 4}, false},
		{Version{1, 2, 3, 5}, Version{1, 2, 3, 4}, false},
		{Version{0, 2, 3, 4}, Version{1, 2, 3, 4}, false},
		{Version{1, 1, 3, 4}, Version{1, 2, 3, 4}, false},
		{Version{1, 2, 2, 4}, Version{1, 2, 3, 4}, false},
		{Version{1, 2, 3, 3}, Version{1, 2, 3, 4}, false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.v.equalTo(tt.o); got != tt.want {
				t.Errorf("'%s'.equalTo('%s') = %v, want %v", tt.v, tt.o, got, tt.want)
			}
		})
	}
}

func TestVersion_greaterThan(t *testing.T) {
	tests := []struct {
		v    Version
		o    Version
		want bool
	}{
		{Version{0, 0, 0, 0}, Version{0, 0, 0, 0}, false},
		{Version{1, 2, 3, 4}, Version{1, 2, 3, 4}, false},
		{Version{2, 2, 3, 4}, Version{1, 2, 3, 4}, true},
		{Version{1, 3, 3, 4}, Version{1, 2, 3, 4}, true},
		{Version{1, 2, 4, 4}, Version{1, 2, 3, 4}, true},
		{Version{1, 2, 3, 5}, Version{1, 2, 3, 4}, true},
		{Version{0, 2, 3, 4}, Version{1, 2, 3, 4}, false},
		{Version{1, 1, 3, 4}, Version{1, 2, 3, 4}, false},
		{Version{1, 2, 2, 4}, Version{1, 2, 3, 4}, false},
		{Version{1, 2, 3, 3}, Version{1, 2, 3, 4}, false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.v.greaterThan(tt.o); got != tt.want {
				t.Errorf("'%s'.greaterThan('%s') = %v, want %v", tt.v, tt.o, got, tt.want)
			}
		})
	}
}

func TestVersion_lessThan(t *testing.T) {
	tests := []struct {
		v    Version
		o    Version
		want bool
	}{
		{Version{0, 0, 0, 0}, Version{0, 0, 0, 0}, false},
		{Version{1, 2, 3, 4}, Version{1, 2, 3, 4}, false},
		{Version{2, 2, 3, 4}, Version{1, 2, 3, 4}, false},
		{Version{1, 3, 3, 4}, Version{1, 2, 3, 4}, false},
		{Version{1, 2, 4, 4}, Version{1, 2, 3, 4}, false},
		{Version{1, 2, 3, 5}, Version{1, 2, 3, 4}, false},
		{Version{0, 2, 3, 4}, Version{1, 2, 3, 4}, true},
		{Version{1, 1, 3, 4}, Version{1, 2, 3, 4}, true},
		{Version{1, 2, 2, 4}, Version{1, 2, 3, 4}, true},
		{Version{1, 2, 3, 3}, Version{1, 2, 3, 4}, true},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.v.lessThan(tt.o); got != tt.want {
				t.Errorf("'%s'.lessThan('%s') = %v, want %v", tt.v, tt.o, got, tt.want)
			}
		})
	}
}