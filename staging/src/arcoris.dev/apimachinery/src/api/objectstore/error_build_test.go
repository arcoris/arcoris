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

package objectstore

import (
	"errors"
	"testing"
)

func TestErrorForBuildsStructuredError(t *testing.T) {
	err := errorFor(ErrorReasonInvalidKey, validKey(), 1, 2, ErrInvalidKey)

	var storeErr *Error
	if !errors.As(err, &storeErr) {
		t.Fatalf("errors.As failed")
	}
	if storeErr.Reason != ErrorReasonInvalidKey {
		t.Fatalf("Reason = %v; want %v", storeErr.Reason, ErrorReasonInvalidKey)
	}
	if storeErr.Expected != 1 || storeErr.Actual != 2 {
		t.Fatalf("revisions = expected %v actual %v; want 1 and 2", storeErr.Expected, storeErr.Actual)
	}
	requireErrorIs(t, err, ErrInvalidKey)
}
