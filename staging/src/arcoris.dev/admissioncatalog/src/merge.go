// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package admissioncatalog

// Merge composes immutable catalogs into a new immutable catalog.
//
// Duplicate declarations are rejected. Nil catalogs are rejected. Components
// may reference kinds declared in another input catalog because Merge processes
// all reason and kind descriptors before component descriptors.
func Merge(catalogs ...*Catalog) (*Catalog, error) {
	builder := newInitializedBuilder()

	for i, catalog := range catalogs {
		if catalog == nil {
			return nil, NilCatalogError{Operation: "merge", Index: i}
		}
		for _, descriptor := range catalog.Reasons() {
			if err := builder.declareReason(descriptor, "catalogs.reasons"); err != nil {
				return nil, err
			}
		}
		for _, descriptor := range catalog.Kinds() {
			if err := builder.declareKind(descriptor, "catalogs.kinds"); err != nil {
				return nil, err
			}
		}
	}

	for _, catalog := range catalogs {
		for _, descriptor := range catalog.Components() {
			if err := builder.declareComponent(descriptor, "catalogs.components"); err != nil {
				return nil, err
			}
		}
	}

	return builder.Build()
}

// MustMerge returns Merge(catalogs...) or panics when composition is invalid.
//
// It is intended for static catalog assembly and tests.
func MustMerge(catalogs ...*Catalog) *Catalog {
	catalog, err := Merge(catalogs...)
	if err != nil {
		panic(err)
	}
	return catalog
}
