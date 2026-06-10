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

package fieldpath

import (
	"errors"
	"testing"
)

func TestMapKey(t *testing.T) {
	key, err := NewMapKey("app")

	requireNoError(t, err)
	requireEqual(t, key.String(), "app")
	requireEqual(t, key.IsZero(), false)
}

func TestMapKeyAcceptsPathLikeTextAsOpaqueKey(t *testing.T) {
	key, err := NewMapKey("app.kubernetes.io/name")

	requireNoError(t, err)
	requireEqual(t, key.String(), "app.kubernetes.io/name")
}

func TestMapKeyRejectsEmptyKey(t *testing.T) {
	_, err := NewMapKey("")

	requireErrorIs(t, err, ErrEmptyMapKey)
}

func TestMapKeyValidateStructureReportsReason(t *testing.T) {
	err := MapKey("").ValidateStructure()

	var pathErr *Error
	if !errors.As(err, &pathErr) {
		t.Fatalf("expected *Error, got %T", err)
	}

	requireEqual(t, pathErr.Reason, ErrorReasonEmptyMapKey)
}

func TestMustMapKeyPanicsOnEmptyKey(t *testing.T) {
	requirePanic(t, func() {
		MustMapKey("")
	})
}
