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
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestCompareStartsAtRoot(t *testing.T) {
	got, err := Compare(value.StringValue("old"), value.StringValue("new"), types.String().Type(), Options{})
	requireNoError(t, err)

	requireResult(t, got, nil, nil, paths(fieldpath.RootPath()))
}

func TestCompareAtPreservesBasePath(t *testing.T) {
	path := rootField("desired", "name")

	got, err := CompareAt(path, value.StringValue("old"), value.StringValue("new"), types.String().Type(), Options{})
	requireNoError(t, err)

	requireResult(t, got, nil, nil, paths(path))
}

func TestCompareAtInvalidBasePathReturnsInvalidDescriptor(t *testing.T) {
	path := fieldpath.RootPath().Index(-1)

	_, err := CompareAt(path, value.StringValue("old"), value.StringValue("new"), types.String().Type(), Options{})

	requireErrorIs(t, err, ErrInvalidDescriptor)
	requireErrorReason(t, err, ErrorReasonInvalidDescriptor)
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
