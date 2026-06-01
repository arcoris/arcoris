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

package meta

import (
	"errors"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

func TestError(t *testing.T) {
	cause := errors.New("nested")
	err := &Error{
		Record: diagnostic.WrapRecord(
			"objectMeta.name",
			ErrInvalidObjectMeta,
			ErrorReasonInvalidForm,
			"name failed",
			cause,
		),
	}

	if !errors.Is(err, ErrInvalidObjectMeta) {
		t.Fatal("Error does not unwrap broad sentinel")
	}
	if !errors.Is(err, cause) {
		t.Fatal("Error does not unwrap nested cause")
	}
	if got := err.Error(); !strings.Contains(got, "objectMeta.name") || !strings.Contains(got, "name failed") {
		t.Fatalf("Error() = %q", got)
	}
}

func TestRootErrorHelpers(t *testing.T) {
	requireErrorIs(
		t,
		invalid("pageToken", ErrInvalidPageToken, ErrorReasonInvalidForm, "bad token"),
		ErrInvalidPageToken,
	)

	requireErrorIs(
		t,
		nested("objectMeta.labels", ErrInvalidObjectMeta, ErrInvalidPageToken),
		ErrInvalidObjectMeta,
	)
}
