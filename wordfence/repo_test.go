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

package wordfence

import (
	"context"
)

type stub struct {
	fixture []Vulnerability
}

func (s stub) fetch(ctx context.Context) (vulnerabilities, error) {
	return s.fixture, nil
}

//func TestRepo_Get(t *testing.T) {
//	tests := []struct {
//		fixture []Vulnerability
//		want    []*wp.Entity
//		wantErr bool
//	}{
//		{
//			fixture: productionExample,
//			want:    productionExampleEntities(),
//			wantErr: false,
//		},
//		{
//			fixture: productionMultiple,
//			want:    productionMultipleEntities(),
//			wantErr: false,
//		},
//		{
//			fixture: scannerExample,
//			want:    scannerExampleEntities(),
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run("", func(t *testing.T) {
//			r := &Repo{
//				Client: stub{
//					fixture: tt.fixture,
//				},
//			}
//			ctx := context.Background()
//
//			got, err := r.Get(ctx)
//
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//
//			diff := cmp.Diff(
//				tt.want,
//				got,
//				cmp.AllowUnexported(wp.Entity{}, version.Range{}, version.Version{}),
//				cmpopts.SortSlices(func(a, b *wp.Entity) bool {
//					return a.Slug() < b.Slug()
//				}),
//				cmpopts.SortSlices(func(a, b *version.Range) bool {
//					return a.String() < b.String()
//				}),
//			)
//			if diff != "" {
//				t.Errorf("Get() mismatch (-want +got):\n%s", diff)
//			}
//		})
//	}
//}

//
//func TestRepo_Get_whereCVENotIn(t *testing.T) {
//	//return []*wp.Entity{core, simpleShippingEdd, wpgsi, wpgsiProfessional, dtChocolate}
//
//	tests := []struct {
//		fixture []Vulnerability
//		cves    []string
//		want    []*wp.Entity
//		wantErr bool
//	}{
//		{
//			fixture: productionMultiple,
//			cves:    nil,
//			want:    productionMultipleEntities(),
//			wantErr: false,
//		},
//		{
//			fixture: productionMultiple,
//			cves:    []string{},
//			want:    productionMultipleEntities(),
//			wantErr: false,
//		},
//		{
//			fixture: productionMultiple,
//			cves:    []string{"CVE-2022-21664"}, // core
//			want:    productionMultipleEntities()[1:],
//			wantErr: false,
//		},
//		{
//			fixture: productionMultiple,
//			cves:    []string{"CVE-2015-9527"}, // simpleShippingEdd
//			want:    append(productionMultipleEntities()[0:1], productionMultipleEntities()[2:]...),
//			wantErr: false,
//		},
//		{
//			fixture: productionMultiple,
//			cves: []string{
//				"CVE-2022-21664", // core
//				"CVE-2015-9527",  // simpleShippingEdd
//			},
//			want:    productionMultipleEntities()[2:],
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run("", func(t *testing.T) {
//			r := &Repo{
//				Client: stub{
//					fixture: tt.fixture,
//				},
//			}
//			ctx := context.Background()
//
//			got, err := r.WhereCVENotIn(tt.cves...).Get(ctx)
//
//			if (err != nil) != tt.wantErr {
//				t.Errorf("WhereCVENotIn().Get() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//
//			diff := cmp.Diff(
//				tt.want,
//				got,
//				cmp.AllowUnexported(wp.Entity{}, version.Range{}, version.Version{}),
//				cmpopts.SortSlices(func(a, b *wp.Entity) bool {
//					return a.Slug() < b.Slug()
//				}),
//				cmpopts.SortSlices(func(a, b *version.Range) bool {
//					return a.String() < b.String()
//				}),
//			)
//			if diff != "" {
//				t.Errorf("WhereCVENotIn().Get() mismatch (-want +got):\n%s", diff)
//			}
//		})
//	}
//}
//
//func TestRepo_Get_whereNotIn(t *testing.T) {
//	//return []*wp.Entity{core, simpleShippingEdd, wpgsi, wpgsiProfessional, dtChocolate}
//
//	tests := []struct {
//		fixture []Vulnerability
//		ids     []string
//		cves    []string
//		want    []*wp.Entity
//		wantErr bool
//	}{
//		{
//			fixture: productionMultiple,
//			ids:     nil,
//			cves:    nil,
//			want:    productionMultipleEntities(),
//			wantErr: false,
//		},
//		{
//			fixture: productionMultiple,
//			ids:     []string{"0114f098-713d-4eef-8643-901f607375de"}, // core
//			cves:    []string{"CVE-2015-9527"},                        // simpleShippingEdd
//			want:    productionMultipleEntities()[2:],
//			wantErr: false,
//		},
//		{
//			fixture: productionMultiple,
//			ids:     []string{"01179ac2-ad68-4a5d-af67-70d57ed611d2"}, // simpleShippingEdd
//			cves:    []string{"CVE-2022-21664"},                       // core
//			want:    productionMultipleEntities()[2:],
//			wantErr: false,
//		},
//		{
//			fixture: productionMultiple,
//			ids:     []string{"01179ac2-ad68-4a5d-af67-70d57ed611d2"}, // simpleShippingEdd
//			cves:    []string{"CVE-2015-9527"},                        // simpleShippingEdd
//			want:    append(productionMultipleEntities()[0:1], productionMultipleEntities()[2:]...),
//			wantErr: false,
//		},
//		{
//			fixture: productionMultiple,
//			ids: []string{
//				"014da588-9494-493e-8659-590b8e8c14a6", // wpgsi & wpgsiProfessional
//				"06fee60a-e96c-49ce-9007-0d402ef46d72", // dtChocolate
//			},
//			cves: []string{
//				"CVE-2022-21664", // core
//				"CVE-2015-9527",  // simpleShippingEdd
//			},
//			want:    nil,
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run("", func(t *testing.T) {
//			r := &Repo{
//				Client: stub{
//					fixture: tt.fixture,
//				},
//			}
//			ctx := context.Background()
//
//			got, err := r.WhereIDNotIn(tt.ids...).WhereCVENotIn(tt.cves...).Get(ctx)
//
//			if (err != nil) != tt.wantErr {
//				t.Errorf("WhereIDNotIn().WhereCVENotIn().Get() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//
//			diff := cmp.Diff(
//				tt.want,
//				got,
//				cmp.AllowUnexported(wp.Entity{}, version.Range{}, version.Version{}),
//				cmpopts.SortSlices(func(a, b *wp.Entity) bool {
//					return a.Slug() < b.Slug()
//				}),
//				cmpopts.SortSlices(func(a, b *version.Range) bool {
//					return a.String() < b.String()
//				}),
//			)
//			if diff != "" {
//				t.Errorf("WhereIDNotIn().WhereCVENotIn().Get() mismatch (-want +got):\n%s", diff)
//			}
//		})
//	}
//}