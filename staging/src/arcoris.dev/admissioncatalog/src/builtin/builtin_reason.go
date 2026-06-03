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
	// ReasonCapacityExhausted reports that bounded live capacity is currently
	// unavailable.
	ReasonCapacityExhausted admission.Reason = "capacity_exhausted"

	// ReasonBudgetExhausted reports that a spend-only budget cannot accept more
	// work.
	ReasonBudgetExhausted admission.Reason = "budget_exhausted"

	// ReasonRateLimited reports that a rate or token gate denied the attempt.
	ReasonRateLimited admission.Reason = "rate_limited"

	// ReasonOverloaded reports that the component is protecting itself from
	// overload.
	ReasonOverloaded admission.Reason = "overloaded"

	// ReasonBackpressured reports that downstream pressure prevented immediate
	// admission.
	ReasonBackpressured admission.Reason = "backpressured"

	// ReasonClosed reports that the component is closed and no longer admits
	// work.
	ReasonClosed admission.Reason = "closed"

	// ReasonDraining reports that the component is intentionally winding down.
	ReasonDraining admission.Reason = "draining"

	// ReasonDeadlineExceeded reports that admission failed because the execution
	// budget was already exhausted.
	ReasonDeadlineExceeded admission.Reason = "deadline_exceeded"

	// ReasonCanceled reports that admission failed because caller-owned
	// cancellation was already observed.
	ReasonCanceled admission.Reason = "canceled"

	// ReasonPolicyDenied reports that a domain policy denied the attempt.
	ReasonPolicyDenied admission.Reason = "policy_denied"
)

// builtinReasonDescriptorLiterals builds the package's built-in reason
// descriptors.
//
// The descriptors are constructed on demand instead of stored in a mutable
// package-level slice. That keeps the built-in catalog copy-safe and reinforces
// that admission has no global registry state.
func reasonDescriptorLiterals() []admissioncatalog.ReasonDescriptor {
	return []admissioncatalog.ReasonDescriptor{
		{
			Reason: admission.ReasonAdmitted,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityEffectNone,
				admissioncatalog.CapabilityEffectCommitted,
				admissioncatalog.CapabilityEffectOwned,
			),
		},
		{
			Reason: admission.ReasonDenied,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Reason: admission.ReasonQueued,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityQueue,
				admissioncatalog.CapabilityEffectQueued,
			),
		},
		{
			Reason: admission.ReasonDeferred,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityDefer,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonCapacityExhausted,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityDefer,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonBudgetExhausted,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityDefer,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonRateLimited,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityDefer,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonOverloaded,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityDefer,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonBackpressured,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityDefer,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonClosed,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonDraining,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityDefer,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonDeadlineExceeded,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonCanceled,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityDefer,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Reason: ReasonPolicyDenied,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectNone,
			),
		},
	}
}

// ReasonDescriptors returns descriptors for the standard reason constants
// provided by this package and the generic core reasons from admission.
//
// The returned slice is a fresh copy. Capabilities are descriptive catalog
// metadata; Result validity remains driven by Decision, Effect, and grant-shape
// invariants.
func ReasonDescriptors() []admissioncatalog.ReasonDescriptor {
	return reasonDescriptorLiterals()
}

// NewReasonRegistry returns a registry populated with standard reasons.
//
// A panic here means the package's own descriptor literals are invalid, which is
// a programming error in admission itself rather than caller input.
func NewReasonRegistry() *admissioncatalog.ReasonRegistry {
	return admissioncatalog.MustReasonRegistry(ReasonDescriptors()...)
}
