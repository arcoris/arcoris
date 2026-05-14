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

func TestInvalidDetailLevelErrorClassifiesWithErrorsIs(t *testing.T) {
	t.Parallel()

	err := InvalidDetailLevelError{Level: DetailLevel(99)}
	if !errors.Is(err, ErrInvalidDetailLevel) {
		t.Fatalf("errors.Is(%v, ErrInvalidDetailLevel) = false, want true", err)
	}
}

func TestInvalidDetailLevelErrorSupportsErrorsAs(t *testing.T) {
	t.Parallel()

	err := error(InvalidDetailLevelError{Level: DetailLevel(99)})

	var levelErr InvalidDetailLevelError
	if !errors.As(err, &levelErr) {
		t.Fatalf("errors.As(%T, InvalidDetailLevelError) = false, want true", err)
	}
	if levelErr.Level != DetailLevel(99) {
		t.Fatalf("InvalidDetailLevelError.Level = %s, want invalid level", levelErr.Level)
	}
}

func TestInvalidDetailLevelErrorMessage(t *testing.T) {
	t.Parallel()

	err := InvalidDetailLevelError{Level: DetailLevel(99)}
	if err.Error() == "" {
		t.Fatal("InvalidDetailLevelError.Error() returned empty message")
	}
}

func TestValidateDetailLevelReturnsInvalidDetailLevelError(t *testing.T) {
	t.Parallel()

	err := validateDetailLevel(DetailLevel(99))
	if !errors.Is(err, ErrInvalidDetailLevel) {
		t.Fatalf("validateDetailLevel(invalid) = %v, want ErrInvalidDetailLevel", err)
	}

	var levelErr InvalidDetailLevelError
	if !errors.As(err, &levelErr) {
		t.Fatalf("validateDetailLevel(invalid) = %T, want InvalidDetailLevelError", err)
	}
	if levelErr.Level != DetailLevel(99) {
		t.Fatalf("InvalidDetailLevelError.Level = %s, want invalid level", levelErr.Level)
	}
}
