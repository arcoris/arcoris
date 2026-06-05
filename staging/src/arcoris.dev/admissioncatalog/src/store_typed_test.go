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

func TestTypedStoresUseDomainKeys(t *testing.T) {
	reasons := newReasonStore()
	reasons.declare(reasonDescriptor(testReason))
	if _, ok := reasons.get(testReason); !ok {
		t.Fatal("reason store did not use reason key")
	}

	kinds := newKindStore()
	kinds.declare(kindDescriptor(testKind))
	if _, ok := kinds.get(testKind); !ok {
		t.Fatal("kind store did not use kind key")
	}

	components := newComponentStore()
	components.declare(componentDescriptor(testComponent, testKind))
	if _, ok := components.get(testComponent); !ok {
		t.Fatal("component store did not use component ID key")
	}
}

func TestTypedStoreInitializersMakeZeroStoresUsable(t *testing.T) {
	var reasons descriptorStore[admission.Reason, ReasonDescriptor]
	initReasonStore(&reasons)
	if !reasons.declare(reasonDescriptor(testReason)) {
		t.Fatal("initialized reason store rejected first descriptor")
	}

	var kinds descriptorStore[admission.ComponentKind, ComponentKindDescriptor]
	initKindStore(&kinds)
	if !kinds.declare(kindDescriptor(testKind)) {
		t.Fatal("initialized kind store rejected first descriptor")
	}

	var components descriptorStore[admission.ComponentID, ComponentDescriptor]
	initComponentStore(&components)
	if !components.declare(componentDescriptor(testComponent, testKind)) {
		t.Fatal("initialized component store rejected first descriptor")
	}
}
