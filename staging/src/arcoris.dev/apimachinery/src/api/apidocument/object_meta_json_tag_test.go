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
	"arcoris.dev/apimachinery/api/meta"
)

func TestObjectMetaJSONTagsMatchAPIDocument(t *testing.T) {
	goType := reflect.TypeOf(meta.ObjectMeta{})
	fields := []struct {
		goName string
		api    apidocument.FieldName
	}{
		{goName: objectMetaGoFieldName, api: apidocument.ObjectMetaFieldName},
		{goName: objectMetaGoFieldNamePrefix, api: apidocument.ObjectMetaFieldGenerateName},
		{goName: objectMetaGoFieldNamespace, api: apidocument.ObjectMetaFieldNamespace},
		{goName: objectMetaGoFieldUID, api: apidocument.ObjectMetaFieldUID},
		{goName: objectMetaGoFieldResourceVersion, api: apidocument.ObjectMetaFieldResourceVersion},
		{goName: objectMetaGoFieldGeneration, api: apidocument.ObjectMetaFieldGeneration},
		{goName: objectMetaGoFieldCreatedAt, api: apidocument.ObjectMetaFieldCreatedAt},
		{goName: objectMetaGoFieldDeletion, api: apidocument.ObjectMetaFieldDeletion},
		{goName: objectMetaGoFieldLabels, api: apidocument.ObjectMetaFieldLabels},
		{goName: objectMetaGoFieldAnnotations, api: apidocument.ObjectMetaFieldAnnotations},
		{goName: objectMetaGoFieldOwnerReferences, api: apidocument.ObjectMetaFieldOwnerReferences},
		{goName: objectMetaGoFieldFinalizers, api: apidocument.ObjectMetaFieldFinalizers},
	}

	for _, field := range fields {
		t.Run(field.goName, func(t *testing.T) {
			assertJSONTagName(t, goType, field.goName, field.api)
			assertJSONTagHasOptions(t, goType, field.goName, jsonTagOptionOmitEmpty)
		})
	}
	assertJSONTagHasOptions(t, goType, objectMetaGoFieldCreatedAt, jsonTagOptionOmitZero)
}

func TestObjectMetaFlattenedJSONFieldsMatchAPIDocument(t *testing.T) {
	goType := reflect.TypeOf(meta.ObjectMeta{})
	got := collectFlattenedJSONFieldNames(t, goType)
	want := []string{
		apidocument.ObjectMetaFieldName.String(),
		apidocument.ObjectMetaFieldGenerateName.String(),
		apidocument.ObjectMetaFieldNamespace.String(),
		apidocument.ObjectMetaFieldUID.String(),
		apidocument.ObjectMetaFieldResourceVersion.String(),
		apidocument.ObjectMetaFieldGeneration.String(),
		apidocument.ObjectMetaFieldCreatedAt.String(),
		apidocument.ObjectMetaFieldDeletion.String(),
		apidocument.ObjectMetaFieldLabels.String(),
		apidocument.ObjectMetaFieldAnnotations.String(),
		apidocument.ObjectMetaFieldOwnerReferences.String(),
		apidocument.ObjectMetaFieldFinalizers.String(),
	}

	if !slices.Equal(got, want) {
		t.Fatalf("object metadata fields = %#v; want %#v", got, want)
	}
}
