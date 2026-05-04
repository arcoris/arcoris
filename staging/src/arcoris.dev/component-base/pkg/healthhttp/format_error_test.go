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

func TestInvalidFormatErrorClassifiesWithErrorsIs(t *testing.T) {
	t.Parallel()

	err := InvalidFormatError{Format: Format(99)}
	if !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("errors.Is(%v, ErrInvalidFormat) = false, want true", err)
	}
}

func TestInvalidFormatErrorSupportsErrorsAs(t *testing.T) {
	t.Parallel()

	err := error(InvalidFormatError{Format: Format(99)})

	var formatErr InvalidFormatError
	if !errors.As(err, &formatErr) {
		t.Fatalf("errors.As(%T, InvalidFormatError) = false, want true", err)
	}
	if formatErr.Format != Format(99) {
		t.Fatalf("InvalidFormatError.Format = %s, want invalid format", formatErr.Format)
	}
}

func TestInvalidFormatErrorMessage(t *testing.T) {
	t.Parallel()

	err := InvalidFormatError{Format: Format(99)}
	if err.Error() == "" {
		t.Fatal("InvalidFormatError.Error() returned empty message")
	}
}

func TestValidateFormatReturnsInvalidFormatError(t *testing.T) {
	t.Parallel()

	err := validateFormat(Format(99))
	if !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("validateFormat(invalid) = %v, want ErrInvalidFormat", err)
	}

	var formatErr InvalidFormatError
	if !errors.As(err, &formatErr) {
		t.Fatalf("validateFormat(invalid) = %T, want InvalidFormatError", err)
	}
	if formatErr.Format != Format(99) {
		t.Fatalf("InvalidFormatError.Format = %s, want invalid format", formatErr.Format)
	}
}
