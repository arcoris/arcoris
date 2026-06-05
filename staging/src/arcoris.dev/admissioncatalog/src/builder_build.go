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

// Build returns an immutable catalog containing the builder's declarations.
//
// The returned catalog is detached from the builder. Mutating the builder after
// Build returns cannot change any previously built catalog.
func (b *Builder) Build() (*Catalog, error) {
	b.requireNonNil()
	b.init()

	for _, descriptor := range b.components.list() {
		if !b.kinds.has(descriptor.Kind) {
			return nil, UnknownComponentKindError{
				ComponentID: descriptor.ID,
				Kind:        descriptor.Kind,
			}
		}
	}

	return &Catalog{
		reasons:    b.reasons.clone(),
		kinds:      b.kinds.clone(),
		components: b.components.clone(),
	}, nil
}

// MustBuild returns b.Build() or panics when the builder contains invalid
// declarations.
func (b *Builder) MustBuild() *Catalog {
	catalog, err := b.Build()
	if err != nil {
		panic(err)
	}
	return catalog
}
