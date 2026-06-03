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

package valuemerge

import (
	"errors"
	"testing"
)

func TestErrorfAtFormatsDetail(t *testing.T) {
	err := errorfAt(root(), ErrInvalidPath, ErrorReasonInvalidPath, "path %s", "bad")

	var mergeError *Error
	if !errors.As(err, &mergeError) {
		t.Fatalf("errors.As(*Error) = false")
	}
	if mergeError.Detail != "path bad" {
		t.Fatalf("detail = %q; want path bad", mergeError.Detail)
	}
}

func TestWrapAtPreservesCause(t *testing.T) {
	cause := errors.New("cause")
	err := wrapAt(root(), ErrInvalidPath, ErrorReasonInvalidPath, "bad path", cause)

	if !errors.Is(err, cause) {
		t.Fatalf("wrapped error does not preserve cause")
	}
}
