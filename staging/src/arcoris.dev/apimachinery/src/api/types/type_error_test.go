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

import (
	"errors"
	"testing"
)

func TestTypeErrorWrapsClassifiedError(t *testing.T) {
	err := typeError("type.field", ErrInvalidType)
	requireEqual(t, errors.Is(err, ErrInvalidType), true)
	requireEqual(t, err.Error(), "types: type.field: invalid type")

	withoutPath := typeError("", ErrInvalidType)
	requireEqual(t, withoutPath.Error(), "types: invalid type")
}

func TestTypeErrorFormatsReasonAndDetail(t *testing.T) {
	err := typeErrorf(
		"type.range",
		ErrInvalidType,
		TypeErrorReasonInvalidRange,
		"minimum must be <= maximum",
	)

	requireEqual(t, errors.Is(err, ErrInvalidType), true)
	requireEqual(t, err.Error(), "types: type.range: invalid type: invalid_range: minimum must be <= maximum")
}

func TestNilTypeErrorMethods(t *testing.T) {
	var err *TypeError
	requireEqual(t, err.Error(), "<nil>")
	requireEqual[error](t, err.Unwrap(), nil)
}
