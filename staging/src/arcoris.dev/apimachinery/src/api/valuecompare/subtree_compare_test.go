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

package valuecompare

import (
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuefieldset"
)

func TestAddSubtreeUsesValueFieldSet(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(types.Field("image").String().Optional()).Descriptor()
	val := value.MustRecordValue(value.MustRecordMember("image", value.StringValue("v1")))

	got, err := newComparer(Options{}).addSubtree(path, val, descriptor, EmptyResult())
	requireNoError(t, err)

	requireResult(t, got, paths(path.Field(testFieldName("image"))), nil, nil)
}

func TestRemoveSubtreeUsesValueFieldSet(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(types.Field("image").String().Optional()).Descriptor()
	val := value.MustRecordValue(value.MustRecordMember("image", value.StringValue("v1")))

	got, err := newComparer(Options{}).removeSubtree(path, val, descriptor, EmptyResult())
	requireNoError(t, err)

	requireResult(t, got, nil, paths(path.Field(testFieldName("image"))), nil)
}

func TestSubtreeExpansionMatchesValueFieldSet(t *testing.T) {
	tests := []struct {
		name       string
		path       fieldpath.Path
		val        value.Value
		descriptor types.Descriptor
	}{
		{
			name:       "scalar",
			path:       rootField("name"),
			val:        value.StringValue("api"),
			descriptor: types.String().Descriptor(),
		},
		{
			name:       "empty record",
			path:       rootField("spec"),
			val:        value.MustRecordValue(),
			descriptor: types.Object(types.Field("image").String().Optional()).Descriptor(),
		},
		{
			name: "nested record",
			path: rootField("template"),
			val: value.MustRecordValue(
				value.MustRecordMember(
					"spec",
					value.MustRecordValue(value.MustRecordMember("image", value.StringValue("v1"))),
				),
			),
			descriptor: types.Object(
				types.Field("spec").Object(types.Field("image").String().Optional()).Optional(),
			).Descriptor(),
		},
		{
			name:       "map",
			path:       rootField("labels"),
			val:        valueRecord("env", "prod"),
			descriptor: types.MapOf(types.String()).Descriptor(),
		},
		{
			name:       "ordered list",
			path:       rootField("args"),
			val:        value.MustListValue(value.StringValue("serve")),
			descriptor: types.ListOf(types.String()).Ordered().Descriptor(),
		},
		{
			name:       "ListMap",
			path:       rootField("conditions"),
			val:        value.MustListValue(conditionValue("Ready", "True")),
			descriptor: conditionsDescriptor(),
		},
		{
			name:       "ListSet",
			path:       rootField("tags"),
			val:        value.MustListValue(value.StringValue("stable"), value.StringValue("fast")),
			descriptor: types.ListOf(types.String()).Set().Descriptor(),
		},
		{
			name:       "unknown preserve opaque",
			path:       rootField("spec"),
			val:        valueRecord("extra", "new"),
			descriptor: types.Object().UnknownFields(types.UnknownPreserveOpaque).Descriptor(),
		},
		{
			name:       "unknown prune",
			path:       rootField("spec"),
			val:        valueRecord("extra", "new"),
			descriptor: types.Object().UnknownFields(types.UnknownPrune).Descriptor(),
		},
	}

	comparer := newComparer(Options{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want, err := valuefieldset.ExtractOwnershipFieldsAt(
				tt.path,
				tt.val,
				tt.descriptor,
				valuefieldset.Options{},
			)
			requireNoError(t, err)

			added, err := comparer.addSubtree(tt.path, tt.val, tt.descriptor, EmptyResult())
			requireNoError(t, err)
			requireSameSet(t, "added subtree", added.Added(), want)
			requireSet(t, "added removed", added.Removed())
			requireSet(t, "added modified", added.Modified())

			removed, err := comparer.removeSubtree(tt.path, tt.val, tt.descriptor, EmptyResult())
			requireNoError(t, err)
			requireSet(t, "removed added", removed.Added())
			requireSameSet(t, "removed subtree", removed.Removed(), want)
			requireSet(t, "removed modified", removed.Modified())
		})
	}
}

func requireSameSet(t *testing.T, name string, got fieldpath.Set, want fieldpath.Set) {
	t.Helper()
	if !got.Equal(want) {
		t.Fatalf("%s = %s, want %s", name, got, want)
	}
}
