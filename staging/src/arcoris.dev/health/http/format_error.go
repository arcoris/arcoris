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
	// ErrInvalidFormat identifies an unsupported health HTTP response format.
	//
	// The sentinel classifies only adapter representation errors. It does not
	// describe Accept negotiation, renderer failures, or report content issues.
	ErrInvalidFormat = errors.New("healthhttp: invalid format")
)

// InvalidFormatError describes an unsupported health HTTP response format.
//
// Callers should use errors.Is for classification and inspect Format only for
// setup diagnostics.
type InvalidFormatError struct {
	Format Format
}

// Error returns the invalid format message.
//
// The exact string is diagnostic-only and may include the format's stable
// textual name.
func (e InvalidFormatError) Error() string {
	return fmt.Sprintf("%v: %s", ErrInvalidFormat, e.Format.String())
}

// Is reports whether target matches the invalid format classification.
//
// This keeps the typed error compatible with errors.Is against ErrInvalidFormat.
func (e InvalidFormatError) Is(target error) bool {
	return target == ErrInvalidFormat
}
