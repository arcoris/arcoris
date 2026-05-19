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

// builtinComponentDescriptorLiterals builds descriptors for known admission
// components already established in the repository architecture.
//
// Component descriptors intentionally represent concrete catalog concepts, not
// every available ComponentKind constant. Additional component packages can add
// owner-created descriptors in their own catalogs later.
func builtinComponentDescriptorLiterals() []ComponentDescriptor {
	return []ComponentDescriptor{
		{
			ID:   "resilience.bulkhead",
			Kind: KindBulkhead,
			Capabilities: NewCapabilitySet(
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityEffectOwned,
				CapabilityEffectNone,
			),
		},
		{
			ID:   "resilience.retrybudget",
			Kind: KindRetryBudget,
			Capabilities: NewCapabilitySet(
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityEffectCommitted,
				CapabilityEffectNone,
			),
		},
		{
			ID:   "resilience.deadline",
			Kind: KindDeadline,
			Capabilities: NewCapabilitySet(
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityDefer,
				CapabilityEffectNone,
			),
		},
	}
}

// BuiltinComponentDescriptors returns descriptors for known built-in admission
// components.
//
// The returned slice is a fresh copy. Component descriptors are stable metadata
// for catalogs and documentation; they are not live component instances.
func BuiltinComponentDescriptors() []ComponentDescriptor {
	return builtinComponentDescriptorLiterals()
}

// NewBuiltinComponentRegistry returns a component registry populated with
// built-in component descriptors.
//
// The caller must pass the KindRegistry used to validate component kinds. A nil
// kind registry panics through MustComponentRegistry, matching the normal
// ComponentRegistry construction policy.
func NewBuiltinComponentRegistry(kinds *KindRegistry) *ComponentRegistry {
	return MustComponentRegistry(kinds, BuiltinComponentDescriptors()...)
}
