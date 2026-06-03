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

package objectapply

import (
	"errors"
	"testing"
)

func TestErrorAt(t *testing.T) {
	err := errorAt(pathObject, ErrInvalidRequest, ErrorReasonInvalidRequest, "bad request")

	requireObjectApplyError(t, err, pathObject, ErrorReasonInvalidRequest)
}

func TestWrapAt(t *testing.T) {
	cause := errors.New("cause")
	err := wrapAt(pathObject, ErrInvalidRequest, ErrorReasonInvalidRequest, "bad request", cause)

	requireObjectApplyError(t, err, pathObject, ErrorReasonInvalidRequest)
	if !errors.Is(err, cause) {
		t.Fatalf("wrapped error does not preserve cause")
	}
}

func TestErrorfAt(t *testing.T) {
	err := errorfAt(pathObject, ErrInvalidRequest, ErrorReasonInvalidRequest, "bad %s", "request")

	requireObjectApplyError(t, err, pathObject, ErrorReasonInvalidRequest)
}
