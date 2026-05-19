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

package admission

// ComponentDescriptor describes stable catalog metadata for an open-world
// admission component.
//
// A descriptor identifies a component role or catalog entry, not a runtime
// instance. It does not register itself globally and it does not imply that an
// instance is currently running. ComponentDescriptor.IsValid performs only local
// syntax checks. ComponentRegistry performs catalog-level checks such as known
// kind membership and duplicate component IDs.
type ComponentDescriptor struct {
	// ID is the stable component identifier or ownership path.
	ID ComponentID

	// Kind is the stable role of the component. Registry-level validation checks
	// that this kind is present in the owner's KindRegistry.
	Kind ComponentKind

	// Capabilities records the outcome and effect classes the component may
	// produce. A zero value is valid and means unspecified. Capabilities are
	// declared behavior surface for catalogs and docs, not enforcement of Result
	// validity.
	Capabilities CapabilitySet
}

// IsValid reports whether d has syntactically valid component metadata.
//
// The method does not require d.Kind to be registered anywhere and it does not
// check duplicate IDs. Those owner-created catalog checks belong to
// ComponentRegistry.
func (d ComponentDescriptor) IsValid() bool {
	return d.ID.IsValid() && d.Kind.IsValid() && d.Capabilities.IsValid()
}
