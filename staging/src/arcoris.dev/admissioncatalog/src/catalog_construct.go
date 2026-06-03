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

// NewCatalog creates an aggregate catalog from owner-provided registries.
//
// All registries are required. Passing nil would create ambiguous partial
// catalog behavior, so construction rejects nil registries explicitly. NewCatalog
// also rejects component registries backed by a different KindRegistry. The
// aggregate catalog assumes a single kind catalog for both kind lookup and
// component validation.
func NewCatalog(
	reasons *ReasonRegistry,
	kinds *KindRegistry,
	components *ComponentRegistry,
) (*Catalog, error) {
	if reasons == nil {
		return nil, ErrNilReasonRegistry
	}
	if kinds == nil {
		return nil, ErrNilKindRegistry
	}
	if components == nil {
		return nil, ErrNilComponentRegistry
	}
	if components.kinds != kinds {
		return nil, ErrMismatchedKindRegistry
	}

	return &Catalog{
		reasons:    reasons,
		kinds:      kinds,
		components: components,
	}, nil
}
