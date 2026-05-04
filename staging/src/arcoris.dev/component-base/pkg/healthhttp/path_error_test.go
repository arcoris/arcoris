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

package healthhttp

import (
	"errors"
	"testing"
)

func TestInvalidPathErrorClassifiesWithErrorsIs(t *testing.T) {
	t.Parallel()

	err := InvalidPathError{Path: "readyz"}
	if !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("errors.Is(%v, ErrInvalidPath) = false, want true", err)
	}
}

func TestInvalidPathErrorSupportsErrorsAs(t *testing.T) {
	t.Parallel()

	err := error(InvalidPathError{Path: "readyz"})

	var pathErr InvalidPathError
	if !errors.As(err, &pathErr) {
		t.Fatalf("errors.As(%T, InvalidPathError) = false, want true", err)
	}
	if pathErr.Path != "readyz" {
		t.Fatalf("InvalidPathError.Path = %q, want readyz", pathErr.Path)
	}
}

func TestInvalidPathErrorMessage(t *testing.T) {
	t.Parallel()

	err := InvalidPathError{Path: "readyz"}
	if err.Error() == "" {
		t.Fatal("InvalidPathError.Error() returned empty message")
	}
}
