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
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/objectsurface"
	"arcoris.dev/apimachinery/api/value"
)

// TestFieldRefStringAndValidate verifies field refs are stable diagnostics,
// not arbitrary string selectors.
func TestFieldRefStringAndValidate(t *testing.T) {
	ref := fieldRef("spec.image")

	if got := ref.String(); got != "desired.spec.image" {
		t.Fatalf("String() = %q; want desired.spec.image", got)
	}
	requireNoError(t, ref.Validate())

	err := FieldRef{Path: fieldPath("$.spec.image")}.Validate()
	requireErrorIs(t, err, ErrInvalidField)

	err = FieldRef{Surface: objectsurface.Kinds().Desired(), Path: fieldpath.Root()}.Validate()
	requireErrorIs(t, err, ErrInvalidField)
}

// TestFieldRefRejectsUnsupportedSurfaces verifies field terms do not silently
// treat metadata or reserved surfaces as missing payload fields.
func TestFieldRefRejectsUnsupportedSurfaces(t *testing.T) {
	metadata := objectsurface.Kinds().Metadata()
	surfaces := []objectsurface.Kind{
		metadata.Labels(),
		metadata.Annotations(),
		metadata.Finalizers(),
		metadata.OwnerReferences(),
	}

	for _, surface := range surfaces {
		t.Run(surface.String(), func(t *testing.T) {
			ref := FieldRef{Surface: surface, Path: fieldPath("$.spec.image")}
			requireErrorIs(t, ref.Validate(), ErrInvalidField)

			_, err := FieldExists(ref)
			requireErrorIs(t, err, ErrInvalidField)
		})
	}
}

// fieldRef builds a desired-surface selectable field reference for tests.
func fieldRef(path string) FieldRef {
	return FieldRef{
		Surface: objectsurface.Kinds().Desired(),
		Path:    fieldPath(path),
	}
}

// observedFieldRef builds an observed-surface selectable field reference for
// tests that need to verify surface dispatch.
func observedFieldRef(path string) FieldRef {
	return FieldRef{
		Surface: objectsurface.Kinds().Observed(),
		Path:    fieldPath(path),
	}
}

// selectable builds a valid selectable field declaration for tests.
func selectable(ref FieldRef, kind value.Kind, operators OperatorSet) SelectableField {
	return SelectableField{
		Ref:       ref,
		Kind:      kind,
		Operators: operators,
		Index:     IndexEquality,
	}
}

// mustFieldSet validates a static field set fixture.
func mustFieldSet(t testing.TB, fields ...SelectableField) StaticFieldSet {
	t.Helper()

	set, err := NewStaticFieldSet(fields...)
	requireNoError(t, err)
	return set
}

// fieldPath parses canonical fieldpath text for selectable field fixtures.
func fieldPath(path string) fieldpath.Path {
	if !strings.HasPrefix(path, "$") {
		path = "$." + path
	}
	parsed, err := fieldpath.ParseCanonical(path)
	if err != nil {
		panic(err)
	}

	return parsed
}
