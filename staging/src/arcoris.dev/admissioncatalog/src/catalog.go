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

import "arcoris.dev/admission"

// Catalog is an immutable read model of admission metadata declarations.
//
// The zero value is a valid empty catalog. Catalog methods are safe for
// concurrent reads after construction. A nil *Catalog receiver panics because
// nil catalog wiring is a programming error, while missing descriptors are
// represented by lookup methods returning false.
type Catalog struct {
	// reasons stores immutable reason declarations by reason.
	reasons descriptorStore[admission.Reason, ReasonDescriptor]

	// kinds stores immutable component kind declarations by kind.
	kinds descriptorStore[admission.ComponentKind, ComponentKindDescriptor]

	// components stores immutable component declarations by component ID.
	components descriptorStore[admission.ComponentID, ComponentDescriptor]
}

// requireNonNil enforces the package-wide nil receiver policy for Catalog.
func (c *Catalog) requireNonNil() {
	if c == nil {
		panic(nilCatalogPanic)
	}
}
