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

package types

import "testing"

func TestMapOfRequiresValidValue(t *testing.T) {
	var expr DescriptorExpr
	desc := MapOf(expr).Descriptor()

	requireErrorIs(t, ValidateLocal(desc), ErrInvalidDescriptor)
}

func TestMapLengthAndDefaultKey(t *testing.T) {
	desc := MapOf(String()).MinEntries(1).MaxEntries(10).Descriptor()

	requireNoError(t, ValidateLocal(desc))
	view, ok := desc.AsMap()
	requireEqual(t, ok, true)
	requireEqual(t, view.Key().Code(), DescriptorString)
	requireEqual(t, view.Value().Code(), DescriptorString)
}

func TestMapKeyDescriptor(t *testing.T) {
	desc := MapOf(String()).Keys(String().MinBytes(1).MaxBytes(253)).Descriptor()

	requireNoError(t, ValidateLocal(desc))
	view, ok := desc.AsMap()
	requireEqual(t, ok, true)

	key, ok := view.Key().AsString()
	requireEqual(t, ok, true)

	minBytes, ok := key.MinBytes()
	requireEqual(t, ok, true)
	requireEqual(t, minBytes, 1)

	maxBytes, ok := key.MaxBytes()
	requireEqual(t, ok, true)
	requireEqual(t, maxBytes, 253)
}

func TestMapInvalidRulesRejected(t *testing.T) {
	invalidLen := MapOf(String()).MinEntries(2).MaxEntries(1).Descriptor()
	invalidKey := MapOf(String()).Keys(Bool()).Descriptor()

	requireErrorIs(t, ValidateLocal(invalidLen), ErrInvalidDescriptor)
	requireErrorIs(t, ValidateLocal(invalidKey), ErrInvalidField)
}

func TestMapDescriptorExprMarker(t *testing.T) {
	MapOf(String()).descriptorExpr()
}
