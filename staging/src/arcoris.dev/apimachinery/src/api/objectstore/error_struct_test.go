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
	"strings"
	"testing"
)

func TestErrorSupportsErrorsIs(t *testing.T) {
	err := &Error{
		Reason:   ErrorReasonStaleRevision,
		Key:      validKey(),
		Expected: 1,
		Actual:   2,
		Err:      ErrStaleRevision,
	}

	requireErrorIs(t, err, ErrStaleRevision)
}

func TestErrorStringIncludesReasonKeyAndRevisions(t *testing.T) {
	err := (&Error{
		Reason:   ErrorReasonStaleRevision,
		Key:      validKey(),
		Expected: 1,
		Actual:   2,
		Err:      ErrStaleRevision,
	}).Error()

	for _, want := range []string{"stale_revision", "control.arcoris.dev/v1:workers/system/main", "expected=1 actual=2"} {
		if !strings.Contains(err, want) {
			t.Fatalf("Error() = %q; want substring %q", err, want)
		}
	}
}

func TestNilErrorIsSafe(t *testing.T) {
	var err *Error

	if got := err.Error(); got != "<nil>" {
		t.Fatalf("Error() = %q; want %q", got, "<nil>")
	}
	if got := err.Unwrap(); got != nil {
		t.Fatalf("Unwrap() = %v; want nil", got)
	}
}
