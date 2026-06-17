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

	"arcoris.dev/apimachinery/api/objectsurface"
	"arcoris.dev/apimachinery/api/value"
)

// TestStaticFieldSetResolveExact verifies field sets resolve only exact
// FieldRef values.
func TestStaticFieldSetResolveExact(t *testing.T) {
	ref := fieldRef("spec.image")
	set := mustFieldSet(t, selectable(ref, value.KindString, Operators(OperatorEquals)))

	if _, ok := set.ResolveSelectableField(ref); !ok {
		t.Fatal("ResolveSelectableField registered ref = false; want true")
	}
	if _, ok := set.ResolveSelectableField(fieldRef("spec.phase")); ok {
		t.Fatal("ResolveSelectableField unknown ref = true; want false")
	}
}

// TestStaticFieldSetLaterDuplicateReplacesEarlier documents the current
// generated-field-set friendly duplicate policy.
func TestStaticFieldSetLaterDuplicateReplacesEarlier(t *testing.T) {
	ref := fieldRef("spec.image")
	set := mustFieldSet(t,
		selectable(ref, value.KindString, Operators(OperatorEquals)),
		selectable(ref, value.KindString, Operators(OperatorEquals, OperatorHasPrefix)),
	)

	field, ok := set.ResolveSelectableField(ref)
	if !ok {
		t.Fatal("ResolveSelectableField duplicate ref = false; want true")
	}
	if !field.Operators.Supports(OperatorHasPrefix) {
		t.Fatal("later duplicate did not replace earlier field")
	}
}

// TestStaticFieldSetKeysByFieldRef verifies lookup uses the structured surface
// and canonical fieldpath identity, not a loose diagnostic suffix.
func TestStaticFieldSetKeysByFieldRef(t *testing.T) {
	desired := fieldRef("spec.image")
	observed := FieldRef{Surface: objectsurface.Kinds().Observed(), Path: desired.Path}
	quoted := fieldRef(`"spec.image"`)
	set := mustFieldSet(t,
		selectable(desired, value.KindString, Operators(OperatorEquals)),
		selectable(observed, value.KindString, Operators(OperatorHasPrefix)),
		selectable(quoted, value.KindString, Operators(OperatorContains)),
	)

	gotDesired, ok := set.ResolveSelectableField(desired)
	if !ok || !gotDesired.Operators.Supports(OperatorEquals) {
		t.Fatalf("desired field lookup = %v, %v; want equals field", gotDesired.Operators, ok)
	}
	gotObserved, ok := set.ResolveSelectableField(observed)
	if !ok || !gotObserved.Operators.Supports(OperatorHasPrefix) {
		t.Fatalf("observed field lookup = %v, %v; want prefix field", gotObserved.Operators, ok)
	}
	gotQuoted, ok := set.ResolveSelectableField(quoted)
	if !ok || !gotQuoted.Operators.Supports(OperatorContains) {
		t.Fatalf("quoted field lookup = %v, %v; want contains field", gotQuoted.Operators, ok)
	}
}
