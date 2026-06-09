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
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuefieldset"
	"arcoris.dev/apimachinery/api/valuevalidation"
	"testing"
)

func TestCompareStartsAtRoot(t *testing.T) {
	got, err := Compare(value.StringValue("old"), value.StringValue("new"), types.String().Descriptor(), Options{})
	requireNoError(t, err)

	requireResult(t, got, nil, nil, paths(fieldpath.RootPath()))
}

func TestCompareAtPreservesBasePath(t *testing.T) {
	path := rootField("desired", "name")

	got, err := CompareAt(path, value.StringValue("old"), value.StringValue("new"), types.String().Descriptor(), Options{})
	requireNoError(t, err)

	requireResult(t, got, nil, nil, paths(path))
}

func TestCompareAtInvalidBasePathReturnsInvalidPath(t *testing.T) {
	path := fieldpath.RootPath().Index(-1)

	_, err := CompareAt(path, value.StringValue("old"), value.StringValue("new"), types.String().Descriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidPath)
	requireErrorNotIs(t, err, ErrInvalidDescriptor)
	requireErrorReason(t, err, ErrorReasonInvalidPath)
}

func TestCompareAtMatchesValidationAndFieldSetPathSemantics(t *testing.T) {
	path := rootField("conditions")
	descriptor := conditionsDescriptor()
	oldValue := value.MustListValue(conditionValue("Ready", "False"))
	newValue := value.MustListValue(conditionValue("Ready", "True"))
	selectorPath := path.Select(readySelector())

	requireNoError(t, valuevalidation.ValidateAt(
		path,
		newValue,
		descriptor,
		valuevalidation.Options{},
	))

	set, err := valuefieldset.ExtractAt(
		path,
		newValue,
		descriptor,
		valuefieldset.Options{},
	)
	requireNoError(t, err)
	requireSet(t, "fieldset", set, selectorPath.Field("type"), selectorPath.Field("status"))

	got, err := CompareAt(path, oldValue, newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(selectorPath.Field("status")))
}
func TestCompareDispatchesScalarDescriptor(t *testing.T) {
	got, err := newComparer(Options{}).compare(
		fieldpath.RootPath(),
		valuepresence.Present(value.StringValue("old")),
		valuepresence.Present(value.StringValue("new")),
		types.String().Descriptor(),
		0,
	)
	requireNoError(t, err)

	requireResult(t, got, nil, nil, paths(fieldpath.RootPath()))
}
func TestComparePresenceBothAbsentIsEmpty(t *testing.T) {
	got, done, err := newComparer(Options{}).comparePresence(
		rootField("name"),
		valuepresence.Absent(),
		valuepresence.Absent(),
		types.String().Descriptor(),
	)
	requireNoError(t, err)

	if !done {
		t.Fatalf("done = false")
	}
	requireResult(t, got, nil, nil, nil)
}

func TestComparePresenceAddedUsesSubtree(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(types.Field("image").String().Optional()).Descriptor()
	newValue := value.MustObjectValue(value.ObjectMember("image", value.StringValue("v1")))

	got, done, err := newComparer(Options{}).comparePresence(
		path,
		valuepresence.Absent(),
		valuepresence.Present(newValue),
		descriptor,
	)
	requireNoError(t, err)

	if !done {
		t.Fatalf("done = false")
	}
	requireResult(t, got, paths(path.Field("image")), nil, nil)
}

func TestComparePresenceRemovedUsesSubtree(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(types.Field("image").String().Optional()).Descriptor()
	oldValue := value.MustObjectValue(value.ObjectMember("image", value.StringValue("v1")))

	got, done, err := newComparer(Options{}).comparePresence(
		path,
		valuepresence.Present(oldValue),
		valuepresence.Absent(),
		descriptor,
	)
	requireNoError(t, err)

	if !done {
		t.Fatalf("done = false")
	}
	requireResult(t, got, nil, paths(path.Field("image")), nil)
}

func TestComparePresenceBothPresentContinues(t *testing.T) {
	_, done, err := newComparer(Options{}).comparePresence(
		rootField("name"),
		valuepresence.Present(value.StringValue("old")),
		valuepresence.Present(value.StringValue("new")),
		types.String().Descriptor(),
	)
	requireNoError(t, err)

	if done {
		t.Fatalf("done = true")
	}
}
func TestRequireComparableInputsRejectsZeroValue(t *testing.T) {
	err := requireComparableInputs(fieldpath.RootPath(), value.Value{}, value.StringValue("x"), types.String().Descriptor())

	requireErrorIs(t, err, ErrInvalidValue)
	requireErrorReason(t, err, ErrorReasonInvalidZero)
}

func TestRequireComparableInputsRejectsInvalidDescriptor(t *testing.T) {
	err := requireComparableInputs(fieldpath.RootPath(), value.StringValue("x"), value.StringValue("y"), types.Descriptor{})

	requireErrorIs(t, err, ErrInvalidDescriptor)
	requireErrorReason(t, err, ErrorReasonInvalidDescriptor)
}

func TestRequireKindRejectsMismatch(t *testing.T) {
	err := requireKind(fieldpath.RootPath(), value.StringValue("x"), value.KindBool, types.DescriptorBool)

	requireErrorIs(t, err, ErrKindMismatch)
	requireErrorReason(t, err, ErrorReasonKindMismatch)
}
func TestCompareNullBothNullIsEmpty(t *testing.T) {
	got, err := newComparer(Options{}).compareNull(rootField("name"), value.NullValue(), value.NullValue())
	requireNoError(t, err)

	requireResult(t, got, nil, nil, nil)
}

func TestCompareNullScalarChangeIsModified(t *testing.T) {
	path := rootField("name")

	got, err := newComparer(Options{}).compareNull(path, value.NullValue(), value.StringValue("x"))
	requireNoError(t, err)

	requireResult(t, got, nil, nil, paths(path))
}

func TestCompareNullDescriptorRejectsNonNull(t *testing.T) {
	_, err := newComparer(Options{}).compareNullDescriptor(
		rootField("name"),
		value.StringValue("x"),
		value.NullValue(),
		types.Null().Descriptor(),
	)

	requireErrorIs(t, err, ErrKindMismatch)
}
func TestCompareUsesPresentOperand(t *testing.T) {
	got := valuepresence.Present(value.NullValue())
	val, ok := got.ValueOK()

	if !ok || !val.IsNull() {
		t.Fatalf("Present() = %#v", got)
	}
}

func TestCompareUsesAbsentOperand(t *testing.T) {
	got := valuepresence.Absent()
	val, ok := got.ValueOK()

	if ok || !val.IsZero() {
		t.Fatalf("Absent() = %#v", got)
	}
}
