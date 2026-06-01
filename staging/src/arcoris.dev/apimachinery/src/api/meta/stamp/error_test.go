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

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

func TestError(t *testing.T) {
	cause := errors.New("nested")
	err := &Error{
		Record: diagnostic.WrapRecord(
			"resourceVersion",
			ErrInvalidResourceVersion,
			ErrorReasonInvalidCharacter,
			"bad token",
			cause,
		),
	}

	if !errors.Is(err, ErrInvalidResourceVersion) || !errors.Is(err, cause) {
		t.Fatal("Error does not unwrap expected sentinels")
	}
	if got := err.Error(); !strings.Contains(got, "resourceVersion") || !strings.Contains(got, "bad token") {
		t.Fatalf("Error() = %q", got)
	}
}

func TestErrorHelpers(t *testing.T) {
	requireErrorIs(
		t,
		invalid("timestamp", ErrInvalidTimestamp, ErrorReasonInvalidForm, "bad timestamp"),
		ErrInvalidTimestamp,
	)
}
