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

// Include declares every descriptor from catalog in the builder.
//
// Include is all-or-nothing. It rejects nil catalogs and duplicate
// declarations, and it leaves the builder unchanged when inclusion fails.
func (b *Builder) Include(catalog *Catalog) error {
	b.requireNonNil()
	b.init()
	if catalog == nil {
		return NilCatalogError{Operation: "include", Index: -1, Path: "include"}
	}

	next := Builder{
		reasons:    b.reasons.clone(),
		kinds:      b.kinds.clone(),
		components: b.components.clone(),
	}

	for i, descriptor := range catalog.Reasons() {
		if err := next.declareReason(descriptor, descriptorPath("include.reasons", i)); err != nil {
			return err
		}
	}
	for i, descriptor := range catalog.Kinds() {
		if err := next.declareKind(descriptor, descriptorPath("include.kinds", i)); err != nil {
			return err
		}
	}
	for i, descriptor := range catalog.Components() {
		if err := next.declareComponent(descriptor, descriptorPath("include.components", i)); err != nil {
			return err
		}
	}
	*b = next
	return nil
}
