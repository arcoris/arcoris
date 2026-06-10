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

func TestElementValidateRejectsEmptyFieldName(t *testing.T) {
	err := (Element{kind: ElementField}).ValidateStructure()

	requireErrorIs(t, err, ErrInvalidElement)
	requireErrorIs(t, err, ErrEmptyFieldName)
}

func TestElementValidateRejectsEmptyKey(t *testing.T) {
	err := (Element{kind: ElementKey}).ValidateStructure()

	requireErrorIs(t, err, ErrInvalidElement)
	requireErrorIs(t, err, ErrEmptyMapKey)
}

func TestElementValidateRejectsNegativeIndex(t *testing.T) {
	err := (Element{kind: ElementIndex, index: -1}).ValidateStructure()

	requireErrorIs(t, err, ErrInvalidElement)
	requireErrorIs(t, err, ErrNegativeIndex)
}

func TestElementValidateRejectsInvalidSelector(t *testing.T) {
	err := (Element{kind: ElementSelector}).ValidateStructure()

	requireErrorIs(t, err, ErrInvalidElement)
	requireErrorIs(t, err, ErrInvalidSelector)
	requireErrorIs(t, err, ErrEmptySelector)
}

func TestElementValidateRejectsInvalidKind(t *testing.T) {
	err := Element{}.ValidateStructure()

	requireErrorIs(t, err, ErrInvalidElement)

	var pathErr *Error
	if !errors.As(err, &pathErr) {
		t.Fatalf("expected *Error, got %T", err)
	}

	requireEqual(t, pathErr.Reason, ErrorReasonInvalidElement)
}
