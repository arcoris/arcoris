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

	"arcoris.dev/apimachinery/api/fieldpath"
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

// TestLookupFieldValueFieldPathElementKinds verifies objectquery traverses
// semantic fieldpath elements instead of splitting diagnostic text.
func TestLookupFieldValueFieldPathElementKinds(t *testing.T) {
	item := testItems()[0]
	tests := []struct {
		name string
		ref  FieldRef
		want string
	}{
		{name: "dynamic key with dots", ref: fieldRef(`$.spec.settings["with.dots"]`), want: "dotted"},
		{name: "list index", ref: fieldRef("$.spec.containers[1].image"), want: "api-sidecar"},
		{name: "associative selector", ref: fieldRef(`$.spec.conditions[{"type":"Ready"}].status`), want: "True"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lookupFieldValue(item, tt.ref)
			if !got.present {
				t.Fatal("field present = false; want true")
			}
			text, ok := got.value.AsString()
			if !ok || text != tt.want {
				t.Fatalf("field value = %q/%v; want %q/true", text, ok, tt.want)
			}
		})
	}
}

// TestLiteralMatchesValue verifies selector matching uses exact literal
// semantics for supported selector scalar kinds.
func TestLiteralMatchesValue(t *testing.T) {
	if !literalMatchesValue(fieldpath.StringLiteral("api"), value.StringValue("api")) {
		t.Fatal("string literal did not match string value")
	}
	if !literalMatchesValue(fieldpath.BoolLiteral(true), value.BoolValue(true)) {
		t.Fatal("bool literal did not match bool value")
	}
	if !literalMatchesValue(fieldpath.Int64Literal(-1), value.Int64Value(-1)) {
		t.Fatal("int64 literal did not match integer value")
	}
	if !literalMatchesValue(fieldpath.Uint64Literal(1), value.Uint64Value(1)) {
		t.Fatal("uint64 literal did not match integer value")
	}
	if literalMatchesValue(fieldpath.StringLiteral("1"), value.Int64Value(1)) {
		t.Fatal("string literal matched integer value")
	}
}
