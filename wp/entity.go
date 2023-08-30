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

package wp

import (
	"errors"

	"github.com/typisttech/wpsecadvi/composer/version"
)

type kind byte

const (
	core kind = iota
	plugin
	theme
)

var (
	ErrEmptyEntitySlug = errors.New("empty entity slug")
)

type Entity struct {
	kind       kind
	slug       string
	constraint version.Constraint
}

func NewCoreEntity() *Entity {
	return &Entity{kind: core, slug: "wordpress-core"}
}

func NewPluginEntity(slug string) (*Entity, error) {
	// TODO: Test me.
	if slug == "" {
		return nil, ErrEmptyEntitySlug
	}
	return &Entity{kind: plugin, slug: slug}, nil
}

func NewThemeEntity(slug string) (*Entity, error) {
	// TODO: Test me.
	if slug == "" {
		return nil, ErrEmptyEntitySlug
	}
	return &Entity{kind: theme, slug: slug}, nil
}

//func (e *Entity) IsCore() bool {
//	return e.kind == core
//}
//
//func (e *Entity) IsPlugin() bool {
//	return e.kind == plugin
//}
//
//func (e *Entity) IsTheme() bool {
//	return e.kind == theme
//}

func (e *Entity) Slug() string {
	return e.slug
}

func (e *Entity) Constraint() version.Constraint {
	return e.constraint
}

func (e *Entity) Or(constraint version.Constraint) {
	e.constraint = append(e.constraint, constraint...)
}