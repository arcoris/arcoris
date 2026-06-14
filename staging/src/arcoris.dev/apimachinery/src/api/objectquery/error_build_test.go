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

package objectquery

import (
	"errors"
	"testing"
)

func TestErrorfBuildsStructuredError(t *testing.T) {
	err := errorf("query.labels", ErrInvalidQuery, ErrorReasonInvalidQuery, "bad %s", "query")

	requireErrorIs(t, err, ErrInvalidQuery)
	var queryErr *Error
	if !errors.As(err, &queryErr) {
		t.Fatalf("errors.As(%T) = false", queryErr)
	}
	if queryErr.Path != "query.labels" {
		t.Fatalf("Path = %q; want query.labels", queryErr.Path)
	}
	if queryErr.Reason != ErrorReasonInvalidQuery {
		t.Fatalf("Reason = %q; want %q", queryErr.Reason, ErrorReasonInvalidQuery)
	}
}

func TestWrapfPreservesCause(t *testing.T) {
	cause := errors.New("lower")
	err := wrapf("query.labels", ErrInvalidQuery, ErrorReasonInvalidQuery, cause, "bad query")

	requireErrorIs(t, err, ErrInvalidQuery)
	requireErrorIs(t, err, cause)
}
