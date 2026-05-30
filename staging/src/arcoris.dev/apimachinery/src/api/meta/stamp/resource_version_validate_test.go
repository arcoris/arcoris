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

package stamp

import (
	"errors"
	"strings"
	"testing"
)

func TestResourceVersionValidate(t *testing.T) {
	requireNoError(t, ResourceVersion("").Validate())
	requireNoError(t, ResourceVersion("rv-1").Validate())

	requireErrorIs(t, ResourceVersion("rv 1").Validate(), ErrInvalidResourceVersion)
	requireErrorIs(t, ResourceVersion("rv/1").Validate(), ErrInvalidResourceVersion)
	requireErrorIs(t, ResourceVersion("rv\n1").Validate(), ErrInvalidResourceVersion)
	requireErrorIs(
		t,
		ResourceVersion(strings.Repeat("x", maxResourceVersionLength+1)).Validate(),
		ErrInvalidResourceVersion,
	)
}

func TestResourceVersionValidateStructuredLengthError(t *testing.T) {
	err := ResourceVersion(strings.Repeat("x", maxResourceVersionLength+1)).Validate()
	requireErrorIs(t, err, ErrInvalidResourceVersion)

	var stampErr *Error
	if !errors.As(err, &stampErr) {
		t.Fatalf("errors.As(%T) = false", stampErr)
	}
	if stampErr.Path != "resourceVersion" {
		t.Fatalf("Path = %q", stampErr.Path)
	}
	if stampErr.Reason != ErrorReasonInvalidLength {
		t.Fatalf("Reason = %q", stampErr.Reason)
	}
	if stampErr.Detail == "" {
		t.Fatal("Detail is empty")
	}
}
