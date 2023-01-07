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

import "fmt"

type PackageType string

const (
	WPCore   PackageType = "wordpress-core"
	WPPlugin PackageType = "wordpress-plugin"
	WPTheme  PackageType = "wordpress-theme"

	WPackagistPluginVendor = "wpackagist-plugin"
	WPackagistThemeVendor  = "wpackagist-theme"
)

var (
	// TODO: Read from yaml files.
	findPackagistOrgCoreFunc SearcherFunc = func(t PackageType, slug string) []string {
		if t != WPCore {
			return nil
		}

		return []string{
			"johnpbloch/wordpress-core",
			"pantheon-systems/wordpress-composer",
			"roots/wordpress-full",
			"roots/wordpress-no-content",
		}
	}
)

type Searcher interface {
	Search(t PackageType, slug string) []string
}

type SearcherFunc func(t PackageType, slug string) []string

func (f SearcherFunc) Search(t PackageType, slug string) []string {
	return f(t, slug)
}

func NewCompositedSearcher() CompositedSearcher {
	cs := CompositedSearcher{}

	cs.AddSearcher(findPackagistOrgCoreFunc)

	return cs
}

type CompositedSearcher struct {
	searchers []Searcher
}

func (cs CompositedSearcher) Search(t PackageType, slug string) []string {
	names := make([]string, 0, len(cs.searchers))
	for _, s := range cs.searchers {
		names = append(
			names,
			s.Search(t, slug)...,
		)
	}

	return names
}

func (cs *CompositedSearcher) AddSearcher(s Searcher) {
	cs.searchers = append(cs.searchers, s)
}

func NewPrefixedSearcher(t PackageType, prefix string) PrefixedSearcher {
	return PrefixedSearcher{
		packageType: t,
		prefix:      prefix,
	}
}

type PrefixedSearcher struct {
	packageType PackageType
	prefix      string
}

func (ps PrefixedSearcher) Search(t PackageType, slug string) []string {
	if t != ps.packageType {
		return nil
	}

	return []string{
		fmt.Sprintf("%s/%s", ps.prefix, slug),
	}
}
