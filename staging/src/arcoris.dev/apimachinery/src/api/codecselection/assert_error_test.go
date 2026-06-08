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

package codecselection

import (
	"errors"
	"strings"
	"testing"
)

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("error = %v; want errors.Is(..., %v)", err, target)
	}
}

func requireSelectionError(t *testing.T, err error, path string, reason ErrorReason) {
	t.Helper()

	var selectionError *Error
	if !errors.As(err, &selectionError) {
		t.Fatalf("error = %T; want *Error", err)
	}
	if selectionError.Path != path {
		t.Fatalf("path = %q; want %q", selectionError.Path, path)
	}
	if selectionError.Reason != reason {
		t.Fatalf("reason = %q; want %q", selectionError.Reason, reason)
	}
}

func requireSelectionDetailContains(t *testing.T, err error, want string) {
	t.Helper()

	var selectionError *Error
	if !errors.As(err, &selectionError) {
		t.Fatalf("error = %T; want *Error", err)
	}
	if !strings.Contains(selectionError.Detail, want) {
		t.Fatalf("detail = %q; want substring %q", selectionError.Detail, want)
	}
}
