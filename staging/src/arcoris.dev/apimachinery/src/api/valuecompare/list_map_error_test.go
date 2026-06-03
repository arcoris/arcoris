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
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/internal/listmapkey"
)

func TestDuplicateListMapEntryErrorIncludesBothIndexes(t *testing.T) {
	err := duplicateListMapEntryError(
		rootField("conditions").Select(readySelector()),
		rootField("conditions").Index(0),
		rootField("conditions").Index(1),
	)

	requireErrorIs(t, err, ErrDuplicateListKey)
	requireErrorDetailContains(t, err, "$.conditions[0]")
	requireErrorDetailContains(t, err, "$.conditions[1]")
}

func TestCompareListMapKeyErrorMapsInternalError(t *testing.T) {
	err := compareListMapKeyError(&listmapkey.Error{
		Path:   rootField("conditions").Index(0).Field("type"),
		Kind:   listmapkey.FailureMissingKey,
		Detail: "missing",
	})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonMissingListKey)
}

func TestCompareListMapKeyErrorLeavesUnknownErrorAlone(t *testing.T) {
	cause := errors.New("plain")

	if got := compareListMapKeyError(cause); !errors.Is(got, cause) {
		t.Fatalf("compareListMapKeyError() did not preserve plain error")
	}
}
