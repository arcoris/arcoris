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

package objectownership

import (
	"errors"
	"testing"
)

func TestErrorAt(t *testing.T) {
	const pathState = "state"

	err := errorAt(pathState, ErrInvalidState, ErrorReasonInvalidState, "bad state")

	requireObjectOwnershipError(t, err, pathState, ErrorReasonInvalidState)
}

func TestErrorfAt(t *testing.T) {
	err := errorfAt(pathState, ErrInvalidState, ErrorReasonInvalidState, "bad %s", "state")

	requireObjectOwnershipError(t, err, pathState, ErrorReasonInvalidState)
}

func TestWrapAt(t *testing.T) {
	const pathState = "state"

	cause := errors.New("cause")
	err := wrapAt(pathState, ErrInvalidState, ErrorReasonInvalidState, "bad state", cause)

	requireObjectOwnershipError(t, err, pathState, ErrorReasonInvalidState)
	requireErrorIs(t, err, cause)
}
