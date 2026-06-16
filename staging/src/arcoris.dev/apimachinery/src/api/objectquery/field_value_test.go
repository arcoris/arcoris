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

package objectquery

import (
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

// TestLookupFieldValueDesiredObservedAndMissing verifies field lookup is
// limited to explicitly supported object surfaces.
func TestLookupFieldValueDesiredObservedAndMissing(t *testing.T) {
	item := testItems()[0]
	observed := desiredRecord("observed-api", "prod", 4)
	item.State.Object.Observed = &observed

	desired := lookupFieldValue(item, fieldRef("spec.image"))
	if !desired.present {
		t.Fatal("desired field present = false; want true")
	}
	got, _ := desired.value.AsString()
	if got != "api" {
		t.Fatalf("desired image = %q; want api", got)
	}

	observedLookup := lookupFieldValue(item, observedFieldRef("spec.image"))
	if !observedLookup.present {
		t.Fatal("observed field present = false; want true")
	}

	missing := lookupFieldValue(item, fieldRef("spec.missing"))
	if missing.present {
		t.Fatal("missing field present = true; want false")
	}

	notRecord := item
	notRecord.State.Object.Desired = value.StringValue("flat")
	if lookupFieldValue(notRecord, fieldRef("spec.image")).present {
		t.Fatal("non-record traversal present = true; want false")
	}
}
