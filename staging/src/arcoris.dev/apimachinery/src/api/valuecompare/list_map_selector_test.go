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

	"arcoris.dev/apimachinery/api/types"
)

func TestExtractListMapSelectorSuccess(t *testing.T) {
	got, err := newComparer(Options{}).extractListMapSelector(
		rootField("conditions").Index(0),
		conditionValue("Ready", "True"),
		conditionDescriptor(),
		[]types.FieldName{"type"},
	)
	requireNoError(t, err)

	if !got.Equal(readySelector()) {
		t.Fatalf("selector = %s, want %s", got, readySelector())
	}
}

func TestExtractListMapSelectorMapsMissingKey(t *testing.T) {
	_, err := newComparer(Options{}).extractListMapSelector(
		rootField("conditions").Index(0),
		valueObject("status", "True"),
		conditionDescriptor(),
		[]types.FieldName{"type"},
	)

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonMissingListKey)
}
