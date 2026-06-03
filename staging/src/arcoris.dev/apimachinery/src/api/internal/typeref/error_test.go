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

package typeref

import (
	"errors"
	"testing"
)

func TestErrorStringIncludesContext(t *testing.T) {
	err := &Error{
		Path:   rootPath(),
		Kind:   FailureUnresolvedRef,
		Detail: "reference was not found",
	}

	got := err.Error()
	if got != "typeref: $: unresolved_ref: reference was not found" {
		t.Fatalf("Error() = %q", got)
	}
}

func TestNilErrorString(t *testing.T) {
	var err *Error
	if got := err.Error(); got != "<nil>" {
		t.Fatalf("Error() = %q; want <nil>", got)
	}
}

func TestErrorUnwrap(t *testing.T) {
	cause := errors.New("cause")
	err := &Error{Cause: cause}

	if got := err.Unwrap(); got != cause {
		t.Fatalf("Unwrap() = %v; want %v", got, cause)
	}
}

func TestNilErrorUnwrap(t *testing.T) {
	var err *Error
	if got := err.Unwrap(); got != nil {
		t.Fatalf("Unwrap() = %v; want nil", got)
	}
}

func TestAsError(t *testing.T) {
	want := &Error{Kind: FailureReferenceCycle}
	got, ok := AsError(want)
	if !ok {
		t.Fatalf("AsError() ok = false")
	}
	if got != want {
		t.Fatalf("AsError() = %p; want %p", got, want)
	}
}

func TestAsErrorRejectsOtherErrors(t *testing.T) {
	got, ok := AsError(errors.New("other"))
	if ok || got != nil {
		t.Fatalf("AsError() = %v, %v; want nil, false", got, ok)
	}
}
