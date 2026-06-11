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
)

func TestRecordFieldPathBuildsCheckedFieldElement(t *testing.T) {
	got, err := recordFieldPath(fieldpath.Root(), "x-y")
	requireNoError(t, err)

	if got.CanonicalText() != `$."x-y"` {
		t.Fatalf("path = %q, want %q", got.CanonicalText(), `$."x-y"`)
	}
}

func TestRecordFieldPathRejectsInvalidName(t *testing.T) {
	_, err := recordFieldPath(fieldpath.Root(), "")

	requireErrorIs(t, err, ErrInvalidPath)
	requireErrorReason(t, err, ErrorReasonInvalidFieldName)
	requireErrorPath(t, err, "$")
}

func TestMapKeyPathBuildsCheckedKeyElement(t *testing.T) {
	got, err := mapKeyPath(fieldpath.Root(), "x-y")
	requireNoError(t, err)

	if got.CanonicalText() != `$["x-y"]` {
		t.Fatalf("path = %q, want %q", got.CanonicalText(), `$["x-y"]`)
	}
}

func TestMapKeyPathRejectsInvalidKey(t *testing.T) {
	_, err := mapKeyPath(fieldpath.Root(), "")

	requireErrorIs(t, err, ErrInvalidPath)
	requireErrorReason(t, err, ErrorReasonInvalidMapKey)
	requireErrorPath(t, err, "$")
}
