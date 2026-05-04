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
	"fmt"
)

var (
	// ErrInvalidHTTPStatusCode identifies an unsupported HTTP status code mapping
	// for health HTTP handlers.
	//
	// Invalid mappings are rejected at configuration boundaries so handlers never
	// emit misleading status codes such as 200 for a failed health target or 404
	// for an internal adapter error.
	ErrInvalidHTTPStatusCode = errors.New("healthhttp: invalid HTTP status code")
)

// InvalidHTTPStatusCodeError describes an invalid HTTP status code field.
//
// InvalidHTTPStatusCodeError is classified as ErrInvalidHTTPStatusCode. Callers
// should use errors.Is for classification and inspect Field and Code only for
// diagnostics.
type InvalidHTTPStatusCodeError struct {
	// Field identifies the mapping field that failed validation.
	//
	// Expected values are "passed", "failed", and "error".
	Field string

	// Code is the invalid configured HTTP status code.
	Code int
}

// Error returns the invalid HTTP status code message.
func (e InvalidHTTPStatusCodeError) Error() string {
	return fmt.Sprintf("%v: field=%s code=%d", ErrInvalidHTTPStatusCode, e.Field, e.Code)
}

// Is reports whether target matches the invalid HTTP status code classification.
func (e InvalidHTTPStatusCodeError) Is(target error) bool {
	return target == ErrInvalidHTTPStatusCode
}
