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

import (
	"testing"

	"arcoris.dev/admission"
)

const (
	testKindBulkhead    admission.ComponentKind = "bulkhead"
	testKindDeadline    admission.ComponentKind = "deadline"
	testKindRetryBudget admission.ComponentKind = "retry_budget"
)

func testReasonRegistry() *ReasonRegistry {
	return MustReasonRegistry(
		testReasonDescriptor(admission.ReasonDenied),
		testReasonDescriptor(admission.ReasonDeferred),
	)
}

func testKindRegistry() *KindRegistry {
	return MustKindRegistry(
		testKindDescriptor(testKindBulkhead),
		testKindDescriptor(testKindDeadline),
		testKindDescriptor(testKindRetryBudget),
	)
}

func testComponentRegistry(kinds *KindRegistry) *ComponentRegistry {
	return MustComponentRegistry(
		kinds,
		testComponentDescriptor("resilience.bulkhead", testKindBulkhead),
		testComponentDescriptor("resilience.deadline", testKindDeadline),
	)
}

func testCatalog() *Catalog {
	kinds := testKindRegistry()
	catalog, err := NewCatalog(
		testReasonRegistry(),
		kinds,
		testComponentRegistry(kinds),
	)
	if err != nil {
		panic(err)
	}
	return catalog
}

func requireCapability(
	t *testing.T,
	set CapabilitySet,
	capability Capability,
) {
	t.Helper()

	if !set.Has(capability) {
		t.Fatalf("capabilities %08b should contain %08b", set, capability)
	}
}
