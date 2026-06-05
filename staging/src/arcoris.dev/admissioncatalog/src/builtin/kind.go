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
	// KindBulkhead describes bounded in-flight isolation components.
	KindBulkhead admission.ComponentKind = "bulkhead"

	// KindRetryBudget describes components that limit retry amplification.
	KindRetryBudget admission.ComponentKind = "retry_budget"

	// KindDeadline describes components that derive or enforce execution
	// budgets from deadlines.
	KindDeadline admission.ComponentKind = "deadline"

	// KindRateLimiter describes components that gate work by rate or token
	// availability.
	KindRateLimiter admission.ComponentKind = "rate_limiter"

	// KindQueue describes components that can take ownership of waiting work.
	KindQueue admission.ComponentKind = "queue"

	// KindScheduler describes components that choose when or where work runs.
	KindScheduler admission.ComponentKind = "scheduler"

	// KindWorkerPool describes components that own worker execution capacity.
	KindWorkerPool admission.ComponentKind = "worker_pool"

	// KindOverloadGate describes components that deny or defer work under
	// overload.
	KindOverloadGate admission.ComponentKind = "overload_gate"

	// KindTenantIsolation describes components that isolate tenants or other
	// workload classes from each other.
	KindTenantIsolation admission.ComponentKind = "tenant_isolation"
)

// KindDescriptors returns fresh descriptors for standard component kinds.
func KindDescriptors() []admissioncatalog.ComponentKindDescriptor {
	return []admissioncatalog.ComponentKindDescriptor{
		{
			Kind:    KindBulkhead,
			Summary: "Bounds live in-flight work and returns caller-owned grants.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityAdmit, admissioncatalog.OutcomeCapabilityDeny),
				effects(admissioncatalog.EffectCapabilityOwned, admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Kind:    KindRetryBudget,
			Summary: "Limits retry amplification through spend-only committed effects.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityAdmit, admissioncatalog.OutcomeCapabilityDeny),
				effects(admissioncatalog.EffectCapabilityCommitted, admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Kind:    KindDeadline,
			Summary: "Derives or checks caller-owned execution budgets.",
			DeclaredCapabilities: capabilities(
				outcomes(
					admissioncatalog.OutcomeCapabilityAdmit,
					admissioncatalog.OutcomeCapabilityDeny,
					admissioncatalog.OutcomeCapabilityDefer,
				),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Kind:    KindRateLimiter,
			Summary: "Gates work by rate or token availability.",
			DeclaredCapabilities: capabilities(
				outcomes(
					admissioncatalog.OutcomeCapabilityAdmit,
					admissioncatalog.OutcomeCapabilityDeny,
					admissioncatalog.OutcomeCapabilityDefer,
				),
				effects(admissioncatalog.EffectCapabilityCommitted, admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Kind:    KindQueue,
			Summary: "Accepts system-owned waiting state.",
			DeclaredCapabilities: capabilities(
				outcomes(
					admissioncatalog.OutcomeCapabilityAdmit,
					admissioncatalog.OutcomeCapabilityDeny,
					admissioncatalog.OutcomeCapabilityQueue,
				),
				effects(admissioncatalog.EffectCapabilityQueued, admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Kind:    KindScheduler,
			Summary: "Chooses when or where admitted work runs.",
			DeclaredCapabilities: capabilities(
				outcomes(
					admissioncatalog.OutcomeCapabilityAdmit,
					admissioncatalog.OutcomeCapabilityDeny,
					admissioncatalog.OutcomeCapabilityDefer,
				),
				effects(admissioncatalog.EffectCapabilityOwned, admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Kind:    KindWorkerPool,
			Summary: "Owns worker execution capacity.",
			DeclaredCapabilities: capabilities(
				outcomes(
					admissioncatalog.OutcomeCapabilityAdmit,
					admissioncatalog.OutcomeCapabilityDeny,
					admissioncatalog.OutcomeCapabilityQueue,
				),
				effects(
					admissioncatalog.EffectCapabilityOwned,
					admissioncatalog.EffectCapabilityQueued,
					admissioncatalog.EffectCapabilityNone,
				),
			),
		},
		{
			Kind:    KindOverloadGate,
			Summary: "Protects a component under overload.",
			DeclaredCapabilities: capabilities(
				outcomes(
					admissioncatalog.OutcomeCapabilityAdmit,
					admissioncatalog.OutcomeCapabilityDeny,
					admissioncatalog.OutcomeCapabilityDefer,
				),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			Kind:    KindTenantIsolation,
			Summary: "Isolates tenants or workload classes.",
			DeclaredCapabilities: capabilities(
				outcomes(
					admissioncatalog.OutcomeCapabilityAdmit,
					admissioncatalog.OutcomeCapabilityDeny,
					admissioncatalog.OutcomeCapabilityDefer,
				),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
	}
}
