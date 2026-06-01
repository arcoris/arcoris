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

func TestErrorIsInvalidElement(t *testing.T) {
	err := FieldElement("").Validate()
	requireErrorIs(t, err, ErrInvalidElement)
}

func TestErrorIsInvalidSelector(t *testing.T) {
	_, err := NewSelector()
	requireErrorIs(t, err, ErrInvalidSelector)
}

func TestErrorIsInvalidLiteral(t *testing.T) {
	err := Literal{}.Validate()
	requireErrorIs(t, err, ErrInvalidLiteral)
}

func TestErrorHasReasonAndDetail(t *testing.T) {
	err := Literal{}.Validate()

	var pathErr *Error
	if !errors.As(err, &pathErr) {
		t.Fatalf("expected *Error, got %T", err)
	}

	requireEqual(t, pathErr.Reason, ErrorReasonInvalidLiteral)
	requireEqual(t, pathErr.Detail != "", true)
}
