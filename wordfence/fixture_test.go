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
	"github.com/typisttech/wpsecadvi/composer/version"
	"github.com/typisttech/wpsecadvi/wp"
)

var (
	productionExample = map[string]Vulnerability{
		"848ccbdc-c6f1-480f-a272-cd459e706713": {
			ID: "848ccbdc-c6f1-480f-a272-cd459e706713",
			Software: []Software{
				{
					Type: "plugin",
					Slug: "example",
					AffectedVersions: map[string]AffectedVersion{
						"1.0.0 - 1.2.3": {
							FromVersion:   "1.0.0",
							FromInclusive: true,
							ToVersion:     "1.2.3",
							ToInclusive:   true,
						},
					},
				},
			},
			CVE: "CVE-1998-1000",
		},
	}

	productionMultiple = map[string]Vulnerability{
		"0114f098-713d-4eef-8643-901f607375de": {
			ID: "0114f098-713d-4eef-8643-901f607375de",
			Software: []Software{
				{
					Type: "core",
					Slug: "wordpress",
					AffectedVersions: map[string]AffectedVersion{
						"[5.6, 5.6.7)": {
							FromVersion:   "5.6",
							FromInclusive: true,
							ToVersion:     "5.6.7",
							ToInclusive:   false,
						},
						"[5.7, 5.7.5)": {
							FromVersion:   "5.7",
							FromInclusive: true,
							ToVersion:     "5.7.5",
							ToInclusive:   false,
						},
						"[5.8, 5.8.3)": {
							FromVersion:   "5.8",
							FromInclusive: true,
							ToVersion:     "5.8.3",
							ToInclusive:   false,
						},
					},
				},
			},
			CVE: "CVE-2022-21664",
		},
		"01179ac2-ad68-4a5d-af67-70d57ed611d2": {
			ID: "01179ac2-ad68-4a5d-af67-70d57ed611d2",
			Software: []Software{
				{
					Type: "plugin",
					Slug: "simple-shipping-edd",
					AffectedVersions: map[string]AffectedVersion{
						"* - 2.1.3": {
							FromVersion:   "*",
							FromInclusive: true,
							ToVersion:     "2.1.3",
							ToInclusive:   true,
						},
					},
				},
			},
			CVE: "CVE-2015-9527",
		},
		"014da588-9494-493e-8659-590b8e8c14a6": {
			ID: "014da588-9494-493e-8659-590b8e8c14a6",
			Software: []Software{
				{
					Type: "plugin",
					Slug: "wpgsi",
					AffectedVersions: map[string]AffectedVersion{
						"* - 3.5.0": {
							FromVersion:   "*",
							FromInclusive: true,
							ToVersion:     "3.5.0",
							ToInclusive:   true,
						},
					},
				},
				{
					Type: "plugin",
					Slug: "wpgsi-professional",
					AffectedVersions: map[string]AffectedVersion{
						"* - 3.5.1": {
							FromVersion:   "*",
							FromInclusive: true,
							ToVersion:     "3.5.1",
							ToInclusive:   true,
						},
					},
				},
			},
		},
		"06fee60a-e96c-49ce-9007-0d402ef46d72": {
			ID: "06fee60a-e96c-49ce-9007-0d402ef46d72",
			Software: []Software{
				{
					Type: "theme",
					Slug: "dt-chocolate",
					AffectedVersions: map[string]AffectedVersion{
						"*": {
							FromVersion:   "*",
							FromInclusive: true,
							ToVersion:     "*",
							ToInclusive:   true,
						},
					},
				},
			},
		},
	}

	scannerExample = map[string]Vulnerability{
		"848ccbdc-c6f1-480f-a272-cd459e706713": {
			ID: "848ccbdc-c6f1-480f-a272-cd459e706713",
			Software: []Software{
				{
					Type: "plugin",
					Slug: "example",
					AffectedVersions: map[string]AffectedVersion{
						"1.0.0 - 1.2.3": {
							FromVersion:   "1.0.0",
							FromInclusive: true,
							ToVersion:     "1.2.3",
							ToInclusive:   true,
						},
					},
				},
			},
		},
	}
)

func productionExampleEntities() []*wp.Entity {
	e, _ := wp.NewPluginEntity("example")

	from, _ := version.New("1.0.0")
	to, _ := version.New("1.2.3")
	r, _ := version.NewRange(
		version.WithInclusiveFloor(from),
		version.WithInclusiveCeiling(to),
	)

	e.Or([]*version.Range{r})

	return []*wp.Entity{e}
}

func productionMultipleEntities() []*wp.Entity {
	core := wp.NewCoreEntity()
	coreFrom1, _ := version.New("5.6")
	coreTo1, _ := version.New("5.6.7")
	coreR1, _ := version.NewRange(
		version.WithInclusiveFloor(coreFrom1),
		version.WithNonInclusiveCeiling(coreTo1),
	)
	coreFrom2, _ := version.New("5.7")
	coreTo2, _ := version.New("5.7.5")
	coreR2, _ := version.NewRange(
		version.WithInclusiveFloor(coreFrom2),
		version.WithNonInclusiveCeiling(coreTo2),
	)
	coreFrom3, _ := version.New("5.8")
	coreTo3, _ := version.New("5.8.3")
	coreR3, _ := version.NewRange(
		version.WithInclusiveFloor(coreFrom3),
		version.WithNonInclusiveCeiling(coreTo3),
	)
	core.Or([]*version.Range{coreR1, coreR2, coreR3})

	simpleShippingEdd, _ := wp.NewPluginEntity("simple-shipping-edd")
	simpleShippingEddTo, _ := version.New("2.1.3")
	simpleShippingEddR, _ := version.NewRange(
		version.WithoutFloor(),
		version.WithInclusiveCeiling(simpleShippingEddTo),
	)
	simpleShippingEdd.Or([]*version.Range{simpleShippingEddR})

	wpgsi, _ := wp.NewPluginEntity("wpgsi")
	wpgsiTo, _ := version.New("3.5.0")
	wpgsiR, _ := version.NewRange(
		version.WithoutFloor(),
		version.WithInclusiveCeiling(wpgsiTo),
	)
	wpgsi.Or([]*version.Range{wpgsiR})

	wpgsiProfessional, _ := wp.NewPluginEntity("wpgsi-professional")
	wpgsiProfessionalTo, _ := version.New("3.5.1")
	wpgsiProfessionalR, _ := version.NewRange(
		version.WithoutFloor(),
		version.WithInclusiveCeiling(wpgsiProfessionalTo),
	)
	wpgsiProfessional.Or([]*version.Range{wpgsiProfessionalR})

	dtChocolate, _ := wp.NewThemeEntity("dt-chocolate")
	dtChocolateR, _ := version.NewRange(
		version.WithoutFloor(),
		version.WithoutCeiling(),
	)
	dtChocolate.Or([]*version.Range{dtChocolateR})

	return []*wp.Entity{core, simpleShippingEdd, wpgsi, wpgsiProfessional, dtChocolate}
}

func scannerExampleEntities() []*wp.Entity {
	return productionExampleEntities()
}