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
	"arcoris.dev/apimachinery/api/types"
	"testing"
)

func TestEqualObjectSameIsTrue(t *testing.T) {
	got, err := newComparer(Options{}).equalObject(rootField("spec"), valueObject("image", "v1"), valueObject("image", "v1"), typesObject("image"), 0)
	requireNoError(t, err)

	if !got {
		t.Fatalf("equalObject() = false")
	}
}

func TestEqualObjectDifferentFieldIsFalse(t *testing.T) {
	got, err := newComparer(Options{}).equalObject(rootField("spec"), valueObject("image", "v1"), valueObject("image", "v2"), typesObject("image"), 0)
	requireNoError(t, err)

	if got {
		t.Fatalf("equalObject() = true")
	}
}

func TestEqualObjectPreservedUnknownSameIsTrue(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserve).Type()

	got, err := newComparer(Options{}).equalObject(rootField("spec"), valueObject("extra", "same"), valueObject("extra", "same"), descriptor, 0)
	requireNoError(t, err)

	if !got {
		t.Fatalf("equalObject() = false")
	}
}

func TestEqualObjectRejectedUnknownReturnsError(t *testing.T) {
	_, err := newComparer(Options{}).equalObject(rootField("spec"), valueObject("extra", "old"), valueObject(), types.Object().Type(), 0)

	requireErrorIs(t, err, ErrUnknownField)
}
