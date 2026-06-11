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

package valuevalidation_test

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valueapply"
	"arcoris.dev/apimachinery/api/valuecompare"
	"arcoris.dev/apimachinery/api/valuefieldset"
	"arcoris.dev/apimachinery/api/valuemerge"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestPipelineListMapSelectorConsistency(t *testing.T) {
	path := rootField("conditions")
	descriptor := conditionListShape()
	oldValue := mustList(t, conditionValue(t, "Ready", "False"))
	newValue := mustList(t, conditionValue(t, "Ready", "True"))
	selectorPath := path.Select(readySelector())
	typePath := selectorPath.Field(fieldpath.MustFieldName("type"))
	statusPath := selectorPath.Field(fieldpath.MustFieldName("status"))

	requireNoError(t, valuevalidation.ValidateAt(path, newValue, descriptor, valuevalidation.Options{}))

	ownershipFields, err := valuefieldset.ExtractOwnershipFieldsAt(path, newValue, descriptor, valuefieldset.Options{})
	requireNoError(t, err)
	requireFieldSet(t, ownershipFields, typePath, statusPath)

	changes, err := valuecompare.CompareAt(path, oldValue, newValue, descriptor, valuecompare.Options{})
	requireNoError(t, err)
	requireFieldSet(t, changes.Changed(), statusPath)

	merged, err := valuemerge.MergeAt(
		path,
		oldValue,
		newValue,
		descriptor,
		fieldpath.MustSet(statusPath),
		valuemerge.Options{},
	)
	requireNoError(t, err)
	requireConditionStatus(t, merged, "True")

	applied, err := valueapply.Apply(valueapply.Request{
		Path:       path,
		Owner:      fieldownership.MustOwner("user"),
		Live:       oldValue,
		Applied:    newValue,
		Descriptor: descriptor,
		Ownership:  fieldownership.EmptyState(),
	}, valueapply.Options{})
	requireNoError(t, err)
	requireFieldSet(t, applied.AppliedFields, typePath, statusPath)
	requireFieldSet(t, applied.ChangedAppliedFields, statusPath)
	requireFieldSet(t, applied.Ownership.FieldsFor(fieldownership.MustOwner("user")), typePath, statusPath)
}

func TestPipelineListSetConservativeOwnershipConsistency(t *testing.T) {
	path := rootField("tags")
	descriptor := types.ListOf(types.String()).Set().Descriptor()
	oldValue := mustList(t, value.StringValue("old"))
	newValue := mustList(t, value.StringValue("new"))

	requireNoError(t, valuevalidation.ValidateAt(path, newValue, descriptor, valuevalidation.Options{}))

	ownershipFields, err := valuefieldset.ExtractOwnershipFieldsAt(path, newValue, descriptor, valuefieldset.Options{})
	requireNoError(t, err)
	requireFieldSet(t, ownershipFields, path)

	changes, err := valuecompare.CompareAt(path, oldValue, newValue, descriptor, valuecompare.Options{})
	requireNoError(t, err)
	requireFieldSet(t, changes.Changed(), path)

	_, err = valuemerge.MergeAt(
		path,
		oldValue,
		newValue,
		descriptor,
		fieldpath.MustSet(path.Index(0)),
		valuemerge.Options{},
	)
	if !errors.Is(err, valuemerge.ErrUnsupportedMerge) {
		t.Fatalf("errors.Is(ErrUnsupportedMerge) = false: %v", err)
	}

	applied, err := valueapply.Apply(valueapply.Request{
		Path:       path,
		Owner:      fieldownership.MustOwner("user"),
		Live:       oldValue,
		Applied:    newValue,
		Descriptor: descriptor,
		Ownership:  fieldownership.EmptyState(),
	}, valueapply.Options{})
	requireNoError(t, err)
	requireFieldSet(t, applied.AppliedFields, path)
	requireFieldSet(t, applied.ChangedAppliedFields, path)
	requireFieldSet(t, applied.Ownership.FieldsFor(fieldownership.MustOwner("user")), path)
}

func TestPipelineUnknownFieldPolicyConsistency(t *testing.T) {
	tests := []struct {
		name          string
		policy        types.UnknownFieldPolicy
		wantFieldSet  []fieldpath.Path
		wantChanged   []fieldpath.Path
		wantErrReject bool
	}{
		{
			name:          "reject",
			policy:        types.UnknownReject,
			wantErrReject: true,
		},
		{
			name:         "preserve opaque",
			policy:       types.UnknownPreserveOpaque,
			wantFieldSet: []fieldpath.Path{rootField("extra")},
			wantChanged:  []fieldpath.Path{rootField("extra")},
		},
		{
			name:   "prune",
			policy: types.UnknownPrune,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			descriptor := types.Object().UnknownFields(tt.policy).Descriptor()
			oldValue := mustObject(t, value.MustRecordMember("extra", value.StringValue("old")))
			newValue := mustObject(t, value.MustRecordMember("extra", value.StringValue("new")))

			validationErr := valuevalidation.Validate(newValue, descriptor, valuevalidation.Options{})
			fieldSet, fieldSetErr := valuefieldset.ExtractOwnershipFields(newValue, descriptor, valuefieldset.Options{})
			compareResult, compareErr := valuecompare.Compare(oldValue, newValue, descriptor, valuecompare.Options{})

			if tt.wantErrReject {
				requireError(t, validationErr, valuevalidation.ErrUnknownField, valuevalidation.ErrorReasonUnknownField, "$.extra")
				if !errors.Is(fieldSetErr, valuefieldset.ErrUnknownField) {
					t.Fatalf("errors.Is(valuefieldset.ErrUnknownField) = false: %v", fieldSetErr)
				}
				if !errors.Is(compareErr, valuecompare.ErrUnknownField) {
					t.Fatalf("errors.Is(valuecompare.ErrUnknownField) = false: %v", compareErr)
				}
				return
			}

			requireNoError(t, validationErr)
			requireNoError(t, fieldSetErr)
			requireNoError(t, compareErr)
			requireFieldSet(t, fieldSet, tt.wantFieldSet...)
			requireFieldSet(t, compareResult.Changed(), tt.wantChanged...)
		})
	}
}

func readySelector() fieldpath.Selector {
	return fieldpath.MustSelector(
		fieldpath.NewSelectorEntry(fieldpath.MustFieldName("type"), fieldpath.StringLiteral("Ready")),
	)
}

func requireConditionStatus(t *testing.T, val value.Value, want string) {
	t.Helper()

	listView, ok := val.AsList()
	if !ok {
		t.Fatalf("value is not a list: %s", val.Kind())
	}
	item, ok := listView.At(0)
	if !ok {
		t.Fatalf("list item 0 is absent")
	}
	recordView, ok := item.AsRecord()
	if !ok {
		t.Fatalf("list item 0 is not a record: %s", item.Kind())
	}
	status, ok := recordView.Get(value.MustMemberName("status"))
	if !ok {
		t.Fatalf("status member is absent")
	}
	got, ok := status.AsString()
	if !ok {
		t.Fatalf("status member is not a string: %s", status.Kind())
	}
	if got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
}

func requireFieldSet(t *testing.T, got fieldpath.Set, want ...fieldpath.Path) {
	t.Helper()

	expected := fieldpath.MustSet(want...)
	if !got.Equal(expected) {
		t.Fatalf("field set = %s, want %s", got, expected)
	}
}
