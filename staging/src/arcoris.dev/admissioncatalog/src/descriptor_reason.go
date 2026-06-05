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

// ReasonDescriptor describes metadata for a stable admission reason.
//
// The descriptor is a copy-safe value. It does not close the open-world
// admission.Reason vocabulary, register itself globally, or enforce admission
// Result validity.
type ReasonDescriptor struct {
	// Reason is the stable machine-readable admission reason.
	Reason admission.Reason

	// Summary is optional human-facing documentation for the reason.
	//
	// An empty summary is valid and means unspecified. Summaries must not contain
	// dynamic request data, tenant IDs, request IDs, secrets, timestamps, stack
	// traces, addresses, or runtime instance identifiers.
	Summary string

	// DeclaredCapabilities describes the usual outcome and side-effect surface
	// associated with the reason.
	//
	// The declaration is metadata only. It does not validate or execute admission
	// results.
	DeclaredCapabilities CapabilitySet
}

// IsValid reports whether d is locally valid descriptor metadata.
//
// Validation checks reason syntax, summary shape, and declared capability bits.
// Duplicate declarations are catalog assembly concerns and are not checked here.
func (d ReasonDescriptor) IsValid() bool {
	return d.Reason.IsValid() &&
		validSummary(d.Summary) &&
		d.DeclaredCapabilities.IsValid()
}

// reasonDescriptorKey returns the catalog key for a reason descriptor.
func reasonDescriptorKey(d ReasonDescriptor) admission.Reason {
	return d.Reason
}

// reasonDescriptorLess orders reason descriptors by stable reason string.
func reasonDescriptorLess(a, b ReasonDescriptor) bool {
	return a.Reason.String() < b.Reason.String()
}
