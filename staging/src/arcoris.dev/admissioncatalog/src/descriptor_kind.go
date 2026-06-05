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

// ComponentKindDescriptor describes metadata for an admission component kind.
//
// The descriptor describes a stable component role, not a live runtime
// component instance.
type ComponentKindDescriptor struct {
	// Kind is the stable component role described by this descriptor.
	Kind admission.ComponentKind

	// Summary is optional human-facing documentation for the component kind.
	//
	// Empty means unspecified. Summaries must remain static catalog metadata and
	// must not carry request, tenant, process, address, or secret data.
	Summary string

	// DeclaredCapabilities describes the common outcome and side-effect surface
	// associated with the component kind.
	DeclaredCapabilities CapabilitySet
}

// IsValid reports whether d is locally valid descriptor metadata.
//
// Validation checks kind syntax, summary shape, and declared capability bits.
// Duplicate declarations are catalog assembly concerns and are not checked here.
func (d ComponentKindDescriptor) IsValid() bool {
	return d.Kind.IsValid() &&
		validSummary(d.Summary) &&
		d.DeclaredCapabilities.IsValid()
}

// kindDescriptorKey returns the catalog key for a component kind descriptor.
func kindDescriptorKey(d ComponentKindDescriptor) admission.ComponentKind {
	return d.Kind
}

// kindDescriptorLess orders kind descriptors by stable kind string.
func kindDescriptorLess(a, b ComponentKindDescriptor) bool {
	return a.Kind.String() < b.Kind.String()
}
