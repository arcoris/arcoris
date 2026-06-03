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

// builtinKindDescriptorLiterals builds the package's built-in kind descriptors.
//
// The descriptors are constructed on demand instead of stored in a mutable
// package-level slice. That keeps the built-in catalog copy-safe even for tests
// in this package and reinforces that admission has no global registry state.
func kindDescriptorLiterals() []admissioncatalog.ComponentKindDescriptor {
	return []admissioncatalog.ComponentKindDescriptor{
		{
			Kind: KindBulkhead,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectOwned,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Kind: KindRetryBudget,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectCommitted,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Kind: KindDeadline,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityDefer,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Kind: KindRateLimiter,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityDefer,
				admissioncatalog.CapabilityEffectCommitted,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Kind: KindQueue,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityQueue,
				admissioncatalog.CapabilityEffectQueued,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Kind: KindScheduler,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityDefer,
				admissioncatalog.CapabilityEffectOwned,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Kind: KindWorkerPool,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityQueue,
				admissioncatalog.CapabilityEffectOwned,
				admissioncatalog.CapabilityEffectQueued,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Kind: KindOverloadGate,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityDefer,
				admissioncatalog.CapabilityEffectNone,
			),
		},
		{
			Kind: KindTenantIsolation,
			Capabilities: admissioncatalog.NewCapabilitySet(
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityDefer,
				admissioncatalog.CapabilityEffectNone,
			),
		},
	}
}

// KindDescriptors returns descriptors for the standard component kind
// constants provided by this package.
//
// The returned slice is a fresh copy. Capabilities are descriptive metadata for
// catalogs and documentation; Result validity remains driven by Decision,
// Effect, and grant-shape invariants.
func KindDescriptors() []admissioncatalog.ComponentKindDescriptor {
	return kindDescriptorLiterals()
}

// NewKindRegistry returns a registry populated with standard descriptors.
//
// A panic here means the package's own descriptor literals are invalid, which is
// a programming error in admission itself rather than caller input.
func NewKindRegistry() *admissioncatalog.KindRegistry {
	return admissioncatalog.MustKindRegistry(KindDescriptors()...)
}
