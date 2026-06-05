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
	ComponentResilienceRetryBudget admission.ComponentID = "resilience.retry_budget"

	// ComponentResilienceDeadline identifies the standard resilience deadline
	// admission component metadata entry.
	ComponentResilienceDeadline admission.ComponentID = "resilience.deadline"
)

// ComponentDescriptors returns fresh descriptors for standard ARCORIS admission
// components already established by the repository.
func ComponentDescriptors() []admissioncatalog.ComponentDescriptor {
	return []admissioncatalog.ComponentDescriptor{
		{
			ID:      ComponentResilienceBulkhead,
			Kind:    KindBulkhead,
			Summary: "Standard resilience bulkhead admission metadata.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityAdmit, admissioncatalog.OutcomeCapabilityDeny),
				effects(admissioncatalog.EffectCapabilityOwned, admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			ID:      ComponentResilienceRetryBudget,
			Kind:    KindRetryBudget,
			Summary: "Standard resilience retry-budget admission metadata.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityAdmit, admissioncatalog.OutcomeCapabilityDeny),
				effects(admissioncatalog.EffectCapabilityCommitted, admissioncatalog.EffectCapabilityNone),
			),
		},
		{
			ID:      ComponentResilienceDeadline,
			Kind:    KindDeadline,
			Summary: "Standard resilience deadline admission metadata.",
			DeclaredCapabilities: capabilities(
				outcomes(admissioncatalog.OutcomeCapabilityAdmit, admissioncatalog.OutcomeCapabilityDeny),
				effects(admissioncatalog.EffectCapabilityNone),
			),
		},
	}
}
