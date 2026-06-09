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

package apidocument_test

import (
	"reflect"
	"slices"
	"testing"

	"arcoris.dev/apimachinery/api/apidocument"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/value"
)

func TestObjectEnvelopeJSONTagsMatchAPIDocument(t *testing.T) {
	goType := reflect.TypeOf(object.Object[value.Value, value.Value]{})

	assertJSONTagNameString(t, goType, objectGoFieldTypeMeta, jsonTagNoExplicitName)
	assertJSONTagHasOptions(t, goType, objectGoFieldTypeMeta, jsonTagOptionInline)

	assertJSONTagName(t, goType, objectGoFieldObjectMeta, apidocument.ObjectFieldMetadata)
	assertJSONTagHasOptions(t, goType, objectGoFieldObjectMeta, jsonTagOptionOmitEmpty, jsonTagOptionOmitZero)

	assertJSONTagName(t, goType, objectGoFieldDesired, apidocument.ObjectFieldDesired)
	assertJSONTagLacksOptions(t, goType, objectGoFieldDesired, jsonTagOptionOmitEmpty, jsonTagOptionOmitZero)

	assertJSONTagName(t, goType, objectGoFieldObserved, apidocument.ObjectFieldObserved)
	assertJSONTagHasOptions(t, goType, objectGoFieldObserved, jsonTagOptionOmitEmpty)
	assertJSONTagLacksOptions(t, goType, objectGoFieldObserved, jsonTagOptionOmitZero)
}

func TestObjectEnvelopeFlattenedJSONFieldsMatchAPIDocument(t *testing.T) {
	goType := reflect.TypeOf(object.Object[value.Value, value.Value]{})
	got := collectFlattenedJSONFieldNames(t, goType)
	want := []string{
		apidocument.ObjectFieldAPIVersion.String(),
		apidocument.ObjectFieldKind.String(),
		apidocument.ObjectFieldMetadata.String(),
		apidocument.ObjectFieldDesired.String(),
		apidocument.ObjectFieldObserved.String(),
	}

	if !slices.Equal(got, want) {
		t.Fatalf("flattened object fields = %#v; want %#v", got, want)
	}
}
