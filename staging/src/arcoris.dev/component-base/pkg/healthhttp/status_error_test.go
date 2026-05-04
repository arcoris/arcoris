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
	"net/http"
	"testing"
)

func TestInvalidHTTPStatusCodeErrorClassifiesWithErrorsIs(t *testing.T) {
	t.Parallel()

	err := InvalidHTTPStatusCodeError{
		Field: "failed",
		Code:  http.StatusOK,
	}

	if !errors.Is(err, ErrInvalidHTTPStatusCode) {
		t.Fatalf("errors.Is(%v, ErrInvalidHTTPStatusCode) = false, want true", err)
	}
}

func TestInvalidHTTPStatusCodeErrorSupportsErrorsAs(t *testing.T) {
	t.Parallel()

	err := error(InvalidHTTPStatusCodeError{
		Field: "failed",
		Code:  http.StatusOK,
	})

	var statusErr InvalidHTTPStatusCodeError
	if !errors.As(err, &statusErr) {
		t.Fatalf("errors.As(%T, InvalidHTTPStatusCodeError) = false, want true", err)
	}
	if statusErr.Field != "failed" {
		t.Fatalf("InvalidHTTPStatusCodeError.Field = %q, want failed", statusErr.Field)
	}
	if statusErr.Code != http.StatusOK {
		t.Fatalf("InvalidHTTPStatusCodeError.Code = %d, want %d", statusErr.Code, http.StatusOK)
	}
}

func TestInvalidHTTPStatusCodeErrorMessage(t *testing.T) {
	t.Parallel()

	err := InvalidHTTPStatusCodeError{
		Field: "error",
		Code:  http.StatusBadRequest,
	}

	if err.Error() == "" {
		t.Fatal("InvalidHTTPStatusCodeError.Error() returned empty message")
	}
}

func TestHTTPStatusCodesValidateReturnsInvalidHTTPStatusCodeError(t *testing.T) {
	t.Parallel()

	err := HTTPStatusCodes{
		Passed: DefaultPassedStatus,
		Failed: http.StatusOK,
		Error:  DefaultErrorStatus,
	}.Validate()

	if !errors.Is(err, ErrInvalidHTTPStatusCode) {
		t.Fatalf("Validate() = %v, want ErrInvalidHTTPStatusCode", err)
	}

	var statusErr InvalidHTTPStatusCodeError
	if !errors.As(err, &statusErr) {
		t.Fatalf("Validate() = %T, want InvalidHTTPStatusCodeError", err)
	}
	if statusErr.Field != "failed" {
		t.Fatalf("InvalidHTTPStatusCodeError.Field = %q, want failed", statusErr.Field)
	}
}
