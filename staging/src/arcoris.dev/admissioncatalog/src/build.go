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

// Build validates input and returns an immutable catalog.
//
// Build is all-or-nothing. It rejects invalid descriptors, duplicate
// declarations, and components whose kind is not declared in the same input. A
// returned Catalog does not share mutable slice storage with input.
func Build(input Input) (*Catalog, error) {
	builder := newInitializedBuilder()

	for i, descriptor := range input.Reasons {
		if err := builder.declareReason(descriptor, descriptorPath("input.reasons", i)); err != nil {
			return nil, err
		}
	}
	for i, descriptor := range input.Kinds {
		if err := builder.declareKind(descriptor, descriptorPath("input.kinds", i)); err != nil {
			return nil, err
		}
	}
	for i, descriptor := range input.Components {
		if err := builder.declareComponent(descriptor, descriptorPath("input.components", i)); err != nil {
			return nil, err
		}
	}

	return builder.Build()
}

// MustBuild returns Build(input) or panics when the input is invalid.
//
// It is intended for static descriptor literals and tests where invalid
// metadata is a programming error.
func MustBuild(input Input) *Catalog {
	catalog, err := Build(input)
	if err != nil {
		panic(err)
	}
	return catalog
}
