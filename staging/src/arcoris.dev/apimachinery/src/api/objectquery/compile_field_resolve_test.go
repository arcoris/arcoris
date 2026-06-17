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

// TestResolveFieldRequiresSelectableFieldSet keeps field terms explicit:
// objectquery does not traverse arbitrary desired/observed paths.
func TestResolveFieldRequiresSelectableFieldSet(t *testing.T) {
	_, err := resolveField(fieldRef("spec.image"), compileOptions{})

	requireErrorIs(t, err, ErrInvalidQuery)
	requireErrorIs(t, err, ErrUnresolvedField)
}

// TestResolveAndValidateFieldUsesRegisteredDefinition verifies successful
// resolution also checks operator support and literal kind.
func TestResolveAndValidateFieldUsesRegisteredDefinition(t *testing.T) {
	ref := fieldRef("spec.image")
	fields := mustFieldSet(t, selectable(ref, value.KindString, Operators(OperatorEquals)))
	query := mustQ(FieldEquals(ref, value.StringValue("api")))
	term := query.expr.term

	field, err := resolveAndValidateField(term, compileOptions{fields: fields})
	requireNoError(t, err)

	if field.Ref.Surface != ref.Surface || !field.Ref.Path.Equal(ref.Path) {
		t.Fatalf("resolved ref = %s; want %s", field.Ref.String(), ref.String())
	}
}
