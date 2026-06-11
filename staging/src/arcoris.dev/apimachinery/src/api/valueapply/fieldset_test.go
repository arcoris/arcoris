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

package valueapply

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func TestExtractOwnershipFieldsAppliedFields(t *testing.T) {
	req := specRequest(owner("user"))

	got, err := newApplier(Options{}).extractAppliedFields(req)
	requireNoError(t, err)

	requireSet(t, got, "$.image")
}

func TestExtractAppliedFieldsEmptyObjectOwnsParentPath(t *testing.T) {
	req := Request{
		Path:       root(),
		Applied:    obj(),
		Descriptor: types.Object(types.Field("name").String().Optional()).Descriptor(),
	}

	got, err := newApplier(Options{}).extractAppliedFields(req)
	requireNoError(t, err)

	requireSet(t, got, "$")
}

func TestExtractAppliedFieldsEmptyMapOwnsParentPath(t *testing.T) {
	req := Request{
		Path:       root(),
		Applied:    obj(),
		Descriptor: mapDescriptor(),
	}

	got, err := newApplier(Options{}).extractAppliedFields(req)
	requireNoError(t, err)

	requireSet(t, got, "$")
}

func TestExtractAppliedFieldsEmptyListOwnsParentPath(t *testing.T) {
	req := Request{
		Path:       root(),
		Applied:    list(),
		Descriptor: orderedStringListDescriptor(),
	}

	got, err := newApplier(Options{}).extractAppliedFields(req)
	requireNoError(t, err)

	requireSet(t, got, "$")
}

func TestExtractAppliedFieldsListSetOwnsParentPath(t *testing.T) {
	req := Request{
		Path:       root(),
		Applied:    list(str("a"), str("b")),
		Descriptor: types.ListOf(types.String()).Set().Descriptor(),
	}

	got, err := newApplier(Options{}).extractAppliedFields(req)
	requireNoError(t, err)

	requireSet(t, got, "$")
}

func TestExtractAppliedFieldsUnknownPreserveOpaqueOwnsUnknownLeaf(t *testing.T) {
	req := Request{
		Path:    root(),
		Applied: obj(member("extra", obj(member("nested", str("value"))))),
		Descriptor: types.Object().
			UnknownFields(types.UnknownPreserveOpaque).
			Descriptor(),
	}

	got, err := newApplier(Options{}).extractAppliedFields(req)
	requireNoError(t, err)

	requireSet(t, got, "$.extra")
}

func TestExtractAppliedFieldsUnknownPruneOmitsUnknownField(t *testing.T) {
	req := Request{
		Path:       root(),
		Applied:    obj(member("extra", str("value"))),
		Descriptor: types.Object().UnknownFields(types.UnknownPrune).Descriptor(),
	}

	got, err := newApplier(Options{}).extractAppliedFields(req)
	requireNoError(t, err)

	requireSet(t, got)
}

func TestFieldSetOptions(t *testing.T) {
	got := newApplier(Options{MaxDepth: 9}).fieldSetOptions()

	if got.MaxDepth != 9 {
		t.Fatalf("MaxDepth = %d; want 9", got.MaxDepth)
	}
}

func TestApplyValueFieldSetErrorWrapped(t *testing.T) {
	req := Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       str("old"),
		Applied:    str("new"),
		Descriptor: types.Descriptor{},
	}

	_, err := newApplier(Options{}).extractAppliedFields(req)

	requireErrorIs(t, err, ErrInvalidDescriptor)
}
