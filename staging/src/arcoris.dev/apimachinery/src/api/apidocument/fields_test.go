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
	"testing"

	"arcoris.dev/apimachinery/api/apidocument"
)

func TestFieldsTreeMatchesFlatConstants(t *testing.T) {
	fields := apidocument.Fields()

	tests := []struct {
		name string
		got  apidocument.FieldName
		want apidocument.FieldName
	}{
		{name: "typeMeta apiVersion", got: fields.TypeMeta().APIVersion(), want: apidocument.TypeMetaFieldAPIVersion},
		{name: "typeMeta kind", got: fields.TypeMeta().Kind(), want: apidocument.TypeMetaFieldKind},
		{name: "object apiVersion", got: fields.Object().APIVersion(), want: apidocument.ObjectFieldAPIVersion},
		{name: "object kind", got: fields.Object().Kind(), want: apidocument.ObjectFieldKind},
		{name: "object metadata", got: fields.Object().Metadata(), want: apidocument.ObjectFieldMetadata},
		{name: "object desired", got: fields.Object().Desired(), want: apidocument.ObjectFieldDesired},
		{name: "object observed", got: fields.Object().Observed(), want: apidocument.ObjectFieldObserved},
		{name: "objectMeta name", got: fields.ObjectMeta().Name(), want: apidocument.ObjectMetaFieldName},
		{name: "objectMeta generateName", got: fields.ObjectMeta().GenerateName(), want: apidocument.ObjectMetaFieldGenerateName},
		{name: "objectMeta namespace", got: fields.ObjectMeta().Namespace(), want: apidocument.ObjectMetaFieldNamespace},
		{name: "objectMeta uid", got: fields.ObjectMeta().UID(), want: apidocument.ObjectMetaFieldUID},
		{name: "objectMeta resourceVersion", got: fields.ObjectMeta().ResourceVersion(), want: apidocument.ObjectMetaFieldResourceVersion},
		{name: "objectMeta generation", got: fields.ObjectMeta().Generation(), want: apidocument.ObjectMetaFieldGeneration},
		{name: "objectMeta createdAt", got: fields.ObjectMeta().CreatedAt(), want: apidocument.ObjectMetaFieldCreatedAt},
		{name: "objectMeta deletion", got: fields.ObjectMeta().Deletion(), want: apidocument.ObjectMetaFieldDeletion},
		{name: "objectMeta labels", got: fields.ObjectMeta().Labels(), want: apidocument.ObjectMetaFieldLabels},
		{name: "objectMeta annotations", got: fields.ObjectMeta().Annotations(), want: apidocument.ObjectMetaFieldAnnotations},
		{name: "objectMeta ownerReferences", got: fields.ObjectMeta().OwnerReferences(), want: apidocument.ObjectMetaFieldOwnerReferences},
		{name: "objectMeta finalizers", got: fields.ObjectMeta().Finalizers(), want: apidocument.ObjectMetaFieldFinalizers},
		{name: "ownership desired", got: fields.Ownership().Desired(), want: apidocument.OwnershipFieldDesired},
		{name: "ownership observed", got: fields.Ownership().Observed(), want: apidocument.OwnershipFieldObserved},
		{name: "ownership metadata", got: fields.Ownership().MetadataField(), want: apidocument.OwnershipFieldMetadata},
		{name: "ownership labels", got: fields.Ownership().Metadata().Labels(), want: apidocument.OwnershipFieldLabels},
		{name: "ownership annotations", got: fields.Ownership().Metadata().Annotations(), want: apidocument.OwnershipFieldAnnotations},
		{name: "ownership entries", got: fields.Ownership().Surface().Entries(), want: apidocument.OwnershipFieldEntries},
		{name: "ownership owner", got: fields.Ownership().Entry().Owner(), want: apidocument.OwnershipFieldOwner},
		{name: "ownership fields", got: fields.Ownership().Entry().Fields(), want: apidocument.OwnershipFieldFields},
		{name: "pageMeta resourceVersion", got: fields.PageMeta().ResourceVersion(), want: apidocument.PageMetaFieldResourceVersion},
		{name: "pageMeta continue", got: fields.PageMeta().Continue(), want: apidocument.PageMetaFieldContinue},
		{name: "pageMeta remainingItemCount", got: fields.PageMeta().RemainingItemCount(), want: apidocument.PageMetaFieldRemainingItemCount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("field = %q; want %q", tt.got, tt.want)
			}
		})
	}
}

func TestFieldGroupsHaveNoDuplicateFields(t *testing.T) {
	fields := apidocument.Fields()
	groups := map[string][]apidocument.FieldName{
		"typeMeta": {
			fields.TypeMeta().APIVersion(),
			fields.TypeMeta().Kind(),
		},
		"object": {
			fields.Object().APIVersion(),
			fields.Object().Kind(),
			fields.Object().Metadata(),
			fields.Object().Desired(),
			fields.Object().Observed(),
		},
		"objectMeta": {
			fields.ObjectMeta().Name(),
			fields.ObjectMeta().GenerateName(),
			fields.ObjectMeta().Namespace(),
			fields.ObjectMeta().UID(),
			fields.ObjectMeta().ResourceVersion(),
			fields.ObjectMeta().Generation(),
			fields.ObjectMeta().CreatedAt(),
			fields.ObjectMeta().Deletion(),
			fields.ObjectMeta().Labels(),
			fields.ObjectMeta().Annotations(),
			fields.ObjectMeta().OwnerReferences(),
			fields.ObjectMeta().Finalizers(),
		},
		"ownership": {
			fields.Ownership().Desired(),
			fields.Ownership().Observed(),
			fields.Ownership().MetadataField(),
			fields.Ownership().Metadata().Labels(),
			fields.Ownership().Metadata().Annotations(),
			fields.Ownership().Surface().Entries(),
			fields.Ownership().Entry().Owner(),
			fields.Ownership().Entry().Fields(),
		},
		"pageMeta": {
			fields.PageMeta().ResourceVersion(),
			fields.PageMeta().Continue(),
			fields.PageMeta().RemainingItemCount(),
		},
	}

	for group, names := range groups {
		t.Run(group, func(t *testing.T) {
			seen := make(map[apidocument.FieldName]struct{}, len(names))
			for _, name := range names {
				if name.IsZero() {
					t.Fatalf("zero field name")
				}
				if _, ok := seen[name]; ok {
					t.Fatalf("duplicate field %q", name)
				}
				seen[name] = struct{}{}
			}
		})
	}
}
