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

package valuemerge

import (
	"testing"

	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/value"
)

func TestPreserveWithoutOverlayContainerPreservesBaseNull(t *testing.T) {
	got, ok, err := preserveWithoutOverlayContainer(
		root().Field(testFieldName("spec")),
		valuepresence.Present(value.NullValue()),
		valuepresence.Absent(),
		value.KindRecord,
	)
	if err != nil {
		t.Fatalf("preserveWithoutOverlayContainer returned error: %v", err)
	}

	if !ok {
		t.Fatalf("ok = false")
	}
	if !got.Value().IsNull() {
		t.Fatalf("value is not null")
	}
}

func TestPreserveWithoutOverlayContainerReturnsAbsentForAbsentBase(t *testing.T) {
	got, ok, err := preserveWithoutOverlayContainer(
		root().Field(testFieldName("spec")),
		valuepresence.Absent(),
		valuepresence.Absent(),
		value.KindRecord,
	)
	if err != nil {
		t.Fatalf("preserveWithoutOverlayContainer returned error: %v", err)
	}

	if !ok {
		t.Fatalf("ok = false")
	}
	if got.Present() {
		t.Fatalf("got is present")
	}
}

func TestPreserveWithoutOverlayContainerContinuesWhenOverlayHasContainer(t *testing.T) {
	_, ok, err := preserveWithoutOverlayContainer(
		root().Field(testFieldName("spec")),
		valuepresence.Absent(),
		valuepresence.Present(obj()),
		value.KindRecord,
	)
	if err != nil {
		t.Fatalf("preserveWithoutOverlayContainer returned error: %v", err)
	}

	if ok {
		t.Fatalf("ok = true")
	}
}

func TestPreserveWithoutOverlayContainerRejectsWrongKindBase(t *testing.T) {
	_, _, err := preserveWithoutOverlayContainer(
		root().Field(testFieldName("spec")),
		valuepresence.Present(str("invalid")),
		valuepresence.Absent(),
		value.KindRecord,
	)

	requireErrorIs(t, err, ErrKindMismatch)
}
