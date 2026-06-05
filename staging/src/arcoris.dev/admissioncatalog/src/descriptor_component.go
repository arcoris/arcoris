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

// ComponentDescriptor describes metadata for a stable admission component.
//
// The descriptor identifies a catalog concept, not a runtime instance. It does
// not imply that a component is currently running or available.
type ComponentDescriptor struct {
	// ID is the stable component metadata identifier.
	ID admission.ComponentID

	// Kind is the stable component role. Catalog assembly validates that the
	// kind has been declared in the assembled catalog.
	Kind admission.ComponentKind

	// Summary is optional human-facing documentation for the component.
	//
	// Empty means unspecified. Summaries must remain static catalog metadata and
	// must not carry runtime instance identifiers, request data, secrets,
	// addresses, timestamps, or stack traces.
	Summary string

	// DeclaredCapabilities describes the outcome and side-effect surface
	// declared for the component.
	DeclaredCapabilities CapabilitySet
}

// IsValid reports whether d is locally valid descriptor metadata.
//
// Validation checks component ID syntax, kind syntax, summary shape, and
// declared capability bits. It does not check whether the kind is declared in a
// catalog and does not check duplicate component IDs.
func (d ComponentDescriptor) IsValid() bool {
	return d.ID.IsValid() &&
		d.Kind.IsValid() &&
		validSummary(d.Summary) &&
		d.DeclaredCapabilities.IsValid()
}

// componentDescriptorKey returns the catalog key for a component descriptor.
func componentDescriptorKey(d ComponentDescriptor) admission.ComponentID {
	return d.ID
}

// componentDescriptorLess orders component descriptors by stable component ID.
func componentDescriptorLess(a, b ComponentDescriptor) bool {
	return a.ID.String() < b.ID.String()
}
