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
)

func TestComparerAbsentAbsentIsEmpty(t *testing.T) {
	got, err := newComparer(Options{}).compare(
		fieldpath.RootPath(),
		absentOperand(),
		absentOperand(),
		types.String().Type(),
		0,
	)
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestComparerAbsentPresentIsAddedSubtree(t *testing.T) {
	got, err := newComparer(Options{}).compare(
		fieldpath.RootPath(),
		absentOperand(),
		presentOperand(value.StringValue("new")),
		types.String().Type(),
		0,
	)
	requireNoError(t, err)
	requireResult(t, got, paths(fieldpath.RootPath()), nil, nil)
}

func TestComparerPresentAbsentIsRemovedSubtree(t *testing.T) {
	got, err := newComparer(Options{}).compare(
		fieldpath.RootPath(),
		presentOperand(value.StringValue("old")),
		absentOperand(),
		types.String().Type(),
		0,
	)
	requireNoError(t, err)
	requireResult(t, got, nil, paths(fieldpath.RootPath()), nil)
}

func TestCompareInvalidZeroValue(t *testing.T) {
	_, err := Compare(value.Value{}, value.StringValue("new"), types.String().Type(), Options{})

	requireErrorIs(t, err, ErrInvalidValue)
	requireErrorReason(t, err, ErrorReasonInvalidZero)
	requireErrorPath(t, err, "$")
}

func TestCompareKindMismatch(t *testing.T) {
	_, err := Compare(value.BoolValue(true), value.StringValue("new"), types.String().Type(), Options{})

	requireErrorIs(t, err, ErrKindMismatch)
	requireErrorReason(t, err, ErrorReasonKindMismatch)
	requireErrorPath(t, err, "$")
}
