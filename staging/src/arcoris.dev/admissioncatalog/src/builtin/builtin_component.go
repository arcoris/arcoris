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

package builtin

import (
	"arcoris.dev/admission"
	"arcoris.dev/admissioncatalog"
)

const (
	// ComponentResilienceBulkhead identifies the standard resilience bulkhead
	// admission component metadata entry.
	ComponentResilienceBulkhead admission.ComponentID = "resilience.bulkhead"

	// ComponentResilienceRetryBudget identifies the standard resilience retry
	// budget admission component metadata entry.
	ComponentResilienceRetryBudget admission.ComponentID = "resilience.retrybudget"

	// ComponentResilienceDeadline identifies the standard resilience deadline
	// admission component metadata entry.
	ComponentResilienceDeadline admission.ComponentID = "resilience.deadline"
)

// builtinComponentDescriptorLiterals builds descriptors for known admission
// components already established in the repository architecture.
//
// Component descriptors intentionally represent concrete catalog concepts, not
// every available ComponentKind constant. Additional component packages can add
// owner-created descriptors in their own catalogs later.
func componentDescriptorLiterals() []admissioncatalog.ComponentDescriptor {
	return []admissioncatalog.ComponentDescriptor{
		{
			ID:   ComponentResilienceBulkhead,
			Kind: KindBulkhead,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectOwned,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			ID:   ComponentResilienceRetryBudget,
			Kind: KindRetryBudget,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectCommitted,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			ID:   ComponentResilienceDeadline,
			Kind: KindDeadline,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectNone,
			),
		},
	}
}

// ComponentDescriptors returns descriptors for known standard admission
// components.
//
// The returned slice is a fresh copy. Component descriptors are stable metadata
// for catalogs and documentation; they are not live component instances.
func ComponentDescriptors() []admissioncatalog.ComponentDescriptor {
	return componentDescriptorLiterals()
}

// NewComponentRegistry returns a component registry populated with standard
// component descriptors.
//
// The caller must pass the KindRegistry used to validate component kinds. A nil
// kind registry panics through MustComponentRegistry, matching the normal
// ComponentRegistry construction policy.
func NewComponentRegistry(kinds *admissioncatalog.KindRegistry) *admissioncatalog.ComponentRegistry {
	return admissioncatalog.MustComponentRegistry(kinds, ComponentDescriptors()...)
}
