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

	// ReasonBudgetExhausted reports that a spend-only retry or admission budget
	// cannot accept more work.
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

// ReasonDescriptors returns fresh descriptors for standard admission reasons.
func ReasonDescriptors() []admissioncatalog.ReasonDescriptor {
	return []admissioncatalog.ReasonDescriptor{
		{
			Reason:  admission.ReasonAdmitted,
			Summary: "Work was admitted immediately.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityAdmit),
				effects(
					admissioncatalog.EffectCapabilityNone,
					admissioncatalog.EffectCapabilityCommitted,
					admissioncatalog.EffectCapabilityOwned,
				),
			),
		},
		{
			Reason:  admission.ReasonDenied,
			Summary: "Work was rejected.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityDeny),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Reason:  admission.ReasonQueued,
			Summary: "Work entered system-owned waiting state.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityQueue),
				effects(admissioncatalog.EffectCapabilityQueued),
			),
		},
		{
			Reason:  admission.ReasonDeferred,
			Summary: "Work was left with the caller for a later retry.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityDefer),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Reason:  ReasonCapacityExhausted,
			Summary: "Bounded live capacity is unavailable.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityDeny, admissioncatalog.OutcomeCapabilityDefer),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Reason:  ReasonBudgetExhausted,
			Summary: "A spend-only budget is exhausted.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityDeny, admissioncatalog.OutcomeCapabilityDefer),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Reason:  ReasonRateLimited,
			Summary: "A rate or token gate denied the attempt.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityDeny, admissioncatalog.OutcomeCapabilityDefer),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Reason:  ReasonOverloaded,
			Summary: "A component is protecting itself from overload.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityDeny, admissioncatalog.OutcomeCapabilityDefer),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Reason:  ReasonBackpressured,
			Summary: "Downstream pressure prevented immediate admission.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityDeny, admissioncatalog.OutcomeCapabilityDefer),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Reason:  ReasonClosed,
			Summary: "A component is closed and no longer admits work.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityDeny),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Reason:  ReasonDraining,
			Summary: "A component is intentionally winding down.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityDeny, admissioncatalog.OutcomeCapabilityDefer),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Reason:  ReasonDeadlineExceeded,
			Summary: "The caller-owned execution budget was already exhausted.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityDeny),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Reason:  ReasonCanceled,
			Summary: "Caller-owned cancellation was already observed.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityDeny, admissioncatalog.OutcomeCapabilityDefer),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Reason:  ReasonPolicyDenied,
			Summary: "A domain policy denied the attempt.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityDeny),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
	}
}
