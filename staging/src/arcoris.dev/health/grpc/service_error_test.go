/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package healthgrpc

import (
	"errors"
	"strings"
	"testing"
)

func TestInvalidServiceError(t *testing.T) {
	t.Parallel()

	err := error(InvalidServiceError{Service: " bad", Index: 3, Reason: "service name must be trimmed"})
	if !errors.Is(err, ErrInvalidService) {
		t.Fatalf("errors.Is(%v, ErrInvalidService) = false, want true", err)
	}

	var serviceErr InvalidServiceError
	if !errors.As(err, &serviceErr) {
		t.Fatalf("errors.As(%T, InvalidServiceError) = false, want true", err)
	}
	if serviceErr.Service != " bad" || serviceErr.Index != 3 || serviceErr.Reason == "" {
		t.Fatalf("InvalidServiceError = %+v, want useful context", serviceErr)
	}
	if !strings.Contains(err.Error(), "service=") || strings.Contains(err.Error(), "password") {
		t.Fatalf("Error() = %q, want service context without unrelated internals", err.Error())
	}
}

func TestDuplicateServiceError(t *testing.T) {
	t.Parallel()

	err := error(DuplicateServiceError{Service: "ready", Index: 2, PreviousIndex: 1})
	if !errors.Is(err, ErrDuplicateService) {
		t.Fatalf("errors.Is(%v, ErrDuplicateService) = false, want true", err)
	}

	var duplicateErr DuplicateServiceError
	if !errors.As(err, &duplicateErr) {
		t.Fatalf("errors.As(%T, DuplicateServiceError) = false, want true", err)
	}
	if duplicateErr.Service != "ready" || duplicateErr.Index != 2 || duplicateErr.PreviousIndex != 1 {
		t.Fatalf("DuplicateServiceError = %+v, want ready indexes 2 and 1", duplicateErr)
	}
	if !strings.Contains(err.Error(), "previous_index=1") {
		t.Fatalf("Error() = %q, want previous_index context", err.Error())
	}
}
