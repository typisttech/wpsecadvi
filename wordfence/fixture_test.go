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
	"github.com/typisttech/wpsecadvi/composer"
)

var (
	productionVulnerabilities = vulnerabilities{
		"123e4567-e89b-12d3-a456-426655440000": {
			ID: "123e4567-e89b-12d3-a456-426655440000",
			Software: []Software{
				{
					Type: "plugin",
					Slug: "foo-bar",
					AffectedVersions: map[string]AffectedVersion{
						"* - 1.5.7": {
							FromVersion:   "*",
							FromInclusive: true,
							ToVersion:     "1.5.7",
							ToInclusive:   true,
						},
						"1.6 - 1.6.3": {
							FromVersion:   "1.6",
							FromInclusive: true,
							ToVersion:     "1.6.3",
							ToInclusive:   true,
						},
						"1.7 - 1.7.3.3": {
							FromVersion:   "1.7",
							FromInclusive: true,
							ToVersion:     "1.7.3.3",
							ToInclusive:   true,
						},
					},
				},
			},
			CVE: "CVE-2022-1111",
		},
		"223e4567-e89b-12d3-a456-426655440000": {
			ID: "223e4567-e89b-12d3-a456-426655440000",
			Software: []Software{
				{
					Type: "plugin",
					Slug: "foo-bar",
					AffectedVersions: map[string]AffectedVersion{
						"[1.8 - 1.8.8.8)": {
							FromVersion:   "1.8",
							FromInclusive: true,
							ToVersion:     "1.8.8.8",
							ToInclusive:   false,
						},
						"(1.9 - 1.9.9.9]": {
							FromVersion:   "1.9",
							FromInclusive: false,
							ToVersion:     "1.9.9.9",
							ToInclusive:   true,
						},
					},
				},
			},
			CVE: "CVE-2022-2222",
		},
		"323e4567-e89b-12d3-a456-426655440000": {
			ID: "323e4567-e89b-12d3-a456-426655440000",
			Software: []Software{
				{
					Type: "theme",
					Slug: "foo-bar",
					AffectedVersions: map[string]AffectedVersion{
						"[2.8 - 2.8.8.8)": {
							FromVersion:   "2.8",
							FromInclusive: true,
							ToVersion:     "2.8.8.8",
							ToInclusive:   false,
						},
						"(2.9 - 2.9.9]": {
							FromVersion:   "2.9",
							FromInclusive: false,
							ToVersion:     "2.9.9",
							ToInclusive:   true,
						},
					},
				},
			},
			CVE: "CVE-2022-3333",
		},
		"423e4567-e89b-12d3-a456-426655440000": {
			ID: "423e4567-e89b-12d3-a456-426655440000",
			Software: []Software{
				{
					Type: "core",
					Slug: "wordpress",
					AffectedVersions: map[string]AffectedVersion{
						"[3.8 - 3.8.8.8)": {
							FromVersion:   "3.8",
							FromInclusive: true,
							ToVersion:     "3.8.8.8",
							ToInclusive:   false,
						},
						"(3.9 - 3.9.9]": {
							FromVersion:   "3.9",
							FromInclusive: false,
							ToVersion:     "3.9.9",
							ToInclusive:   true,
						},
					},
				},
			},
			CVE: "CVE-2022-4444",
		},
		"523e4567-e89b-12d3-a456-426655440000": {
			ID: "523e4567-e89b-12d3-a456-426655440000",
			Software: []Software{
				{
					Type: "plugin",
					Slug: "foo-bar-wildcard",
					AffectedVersions: map[string]AffectedVersion{
						"* - *": {
							FromVersion:   "*",
							FromInclusive: true,
							ToVersion:     "*",
							ToInclusive:   true,
						},
					},
				},
			},
			CVE: "CVE-2022-5555",
		},
		"623e4567-e89b-12d3-a456-426655440000": {
			ID: "623e4567-e89b-12d3-a456-426655440000",
			Software: []Software{
				{
					Type: "plugin",
					Slug: "foo-bar",
					AffectedVersions: map[string]AffectedVersion{
						"[1.0 - 2.0]": {
							FromVersion:   "1.0",
							FromInclusive: true,
							ToVersion:     "2.0",
							ToInclusive:   true,
						},
						"[9.1.1 - 10.2.2)": {
							FromVersion:   "9.1.1",
							FromInclusive: true,
							ToVersion:     "10.2.2",
							ToInclusive:   false,
						},
					},
				},
			},
			CVE: "CVE-2022-6666",
		},
		"723e4567-e89b-12d3-a456-426655440000": {
			ID: "723e4567-e89b-12d3-a456-426655440000",
			Software: []Software{
				{
					Type: "plugin",
					Slug: "foo-bar-informational",
					AffectedVersions: map[string]AffectedVersion{
						"informational": {
							FromVersion:   "informational",
							FromInclusive: true,
							ToVersion:     "informational",
							ToInclusive:   true,
						},
					},
				},
			},
			CVE: "CVE-2022-7777",
		},
		"823e4567-e89b-12d3-a456-426655440000": {
			ID: "823e4567-e89b-12d3-a456-426655440000",
			Software: []Software{
				{
					Type: "core",
					Slug: "wordpress",
					AffectedVersions: map[string]AffectedVersion{
						"* - 6.1.1": {
							FromVersion:   "*",
							FromInclusive: true,
							ToVersion:     "6.1.1",
							ToInclusive:   true,
						},
					},
				},
			},
			CVE: "CVE-2022-8888",
		},
	}

	productionSoftwareRanges = map[composer.PackageType]map[softwareSlug][]string{
		composer.WPCore: {
			"wordpress": []string{">3.9.0,<=3.9.9"},
		},
		composer.WPPlugin: {
			"foo-bar-wildcard": []string{"*"},
			"foo-bar":          []string{"<=1.5.7", ">=1.6.0,<=1.6.3", ">=1.0.0,<=2.0.0", ">=9.1.1,<10.2.2"},
		},
		composer.WPTheme: {
			"foo-bar": []string{">2.9.0,<=2.9.9"},
		},
	}
)
