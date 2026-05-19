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

// builtinReasonDescriptorLiterals builds the package's built-in reason
// descriptors.
//
// The descriptors are constructed on demand instead of stored in a mutable
// package-level slice. That keeps the built-in catalog copy-safe and reinforces
// that admission has no global registry state.
func builtinReasonDescriptorLiterals() []ReasonDescriptor {
	return []ReasonDescriptor{
		{
			Reason: ReasonAdmitted,
			Capabilities: NewCapabilitySet(
				CapabilityAdmit,
				CapabilityEffectNone,
				CapabilityEffectCommitted,
				CapabilityEffectOwned,
			),
		},
		{
			Reason: ReasonDenied,
			Capabilities: NewCapabilitySet(
				CapabilityDeny,
				CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonQueued,
			Capabilities: NewCapabilitySet(
				CapabilityQueue,
				CapabilityEffectQueued,
			),
		},
		{
			Reason: ReasonDeferred,
			Capabilities: NewCapabilitySet(
				CapabilityDefer,
				CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonCapacityExhausted,
			Capabilities: NewCapabilitySet(
				CapabilityDeny,
				CapabilityDefer,
				CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonBudgetExhausted,
			Capabilities: NewCapabilitySet(
				CapabilityDeny,
				CapabilityDefer,
				CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonRateLimited,
			Capabilities: NewCapabilitySet(
				CapabilityDeny,
				CapabilityDefer,
				CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonOverloaded,
			Capabilities: NewCapabilitySet(
				CapabilityDeny,
				CapabilityDefer,
				CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonBackpressured,
			Capabilities: NewCapabilitySet(
				CapabilityDeny,
				CapabilityDefer,
				CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonClosed,
			Capabilities: NewCapabilitySet(
				CapabilityDeny,
				CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonDraining,
			Capabilities: NewCapabilitySet(
				CapabilityDeny,
				CapabilityDefer,
				CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonDeadlineExceeded,
			Capabilities: NewCapabilitySet(
				CapabilityDeny,
				CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonCanceled,
			Capabilities: NewCapabilitySet(
				CapabilityDeny,
				CapabilityDefer,
				CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonPolicyDenied,
			Capabilities: NewCapabilitySet(
				CapabilityDeny,
				CapabilityEffectNone,
			),
		},
	}
}

// BuiltinReasonDescriptors returns descriptors for the reason constants
// provided by this package.
//
// The returned slice is a fresh copy. Capabilities are descriptive catalog
// metadata; Result validity remains driven by Decision, Effect, and grant-shape
// invariants.
func BuiltinReasonDescriptors() []ReasonDescriptor {
	return builtinReasonDescriptorLiterals()
}

// NewBuiltinReasonRegistry returns a registry populated with built-in reasons.
//
// A panic here means the package's own descriptor literals are invalid, which is
// a programming error in admission itself rather than caller input.
func NewBuiltinReasonRegistry() *ReasonRegistry {
	return MustReasonRegistry(BuiltinReasonDescriptors()...)
}
