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

// builtinKindDescriptorLiterals builds the package's built-in kind descriptors.
//
// The descriptors are constructed on demand instead of stored in a mutable
// package-level slice. That keeps the built-in catalog copy-safe even for tests
// in this package and reinforces that admission has no global registry state.
func builtinKindDescriptorLiterals() []ComponentKindDescriptor {
	return []ComponentKindDescriptor{
		{
			Kind: KindBulkhead,
			Capabilities: NewCapabilitySet(
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityEffectOwned,
				CapabilityEffectNone,
			),
		},
		{
			Kind: KindRetryBudget,
			Capabilities: NewCapabilitySet(
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityEffectCommitted,
				CapabilityEffectNone,
			),
		},
		{
			Kind: KindDeadline,
			Capabilities: NewCapabilitySet(
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityDefer,
				CapabilityEffectNone,
			),
		},
		{
			Kind: KindRateLimiter,
			Capabilities: NewCapabilitySet(
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityDefer,
				CapabilityEffectCommitted,
				CapabilityEffectNone,
			),
		},
		{
			Kind: KindQueue,
			Capabilities: NewCapabilitySet(
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityQueue,
				CapabilityEffectQueued,
				CapabilityEffectNone,
			),
		},
		{
			Kind: KindScheduler,
			Capabilities: NewCapabilitySet(
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityDefer,
				CapabilityEffectOwned,
				CapabilityEffectNone,
			),
		},
		{
			Kind: KindWorkerPool,
			Capabilities: NewCapabilitySet(
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityQueue,
				CapabilityEffectOwned,
				CapabilityEffectQueued,
				CapabilityEffectNone,
			),
		},
		{
			Kind: KindOverloadGate,
			Capabilities: NewCapabilitySet(
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityDefer,
				CapabilityEffectNone,
			),
		},
		{
			Kind: KindTenantIsolation,
			Capabilities: NewCapabilitySet(
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityDefer,
				CapabilityEffectNone,
			),
		},
	}
}

// BuiltinKindDescriptors returns descriptors for the kind constants provided by
// this package.
//
// The returned slice is a fresh copy. Capabilities are descriptive metadata for
// catalogs and documentation; Result validity remains driven by Decision,
// Effect, and grant-shape invariants.
func BuiltinKindDescriptors() []ComponentKindDescriptor {
	return builtinKindDescriptorLiterals()
}

// NewBuiltinKindRegistry returns a registry populated with built-in descriptors.
//
// A panic here means the package's own descriptor literals are invalid, which is
// a programming error in admission itself rather than caller input.
func NewBuiltinKindRegistry() *KindRegistry {
	return MustKindRegistry(BuiltinKindDescriptors()...)
}
