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
	// ErrInvalidPath identifies an HTTP path that cannot be used as a health
	// endpoint route.
	//
	// The sentinel covers only adapter-owned path safety invariants. It does not
	// describe router-specific duplicate registration or pattern-matching errors,
	// which remain the responsibility of the concrete mux.
	ErrInvalidPath = errors.New("healthhttp: invalid path")
)

// InvalidPathError describes an invalid health HTTP path.
//
// InvalidPathError is classified as ErrInvalidPath. Callers should use
// errors.Is for classification and inspect Path only for diagnostics. The error
// keeps the original path string so setup failures can be traced without
// widening the public runtime response surface.
type InvalidPathError struct {
	Path string
}

// Error returns the invalid path message.
//
// The formatting is diagnostic-only and is not intended to be a stable public
// protocol.
func (e InvalidPathError) Error() string {
	return fmt.Sprintf("%v: %q", ErrInvalidPath, e.Path)
}

// Is reports whether target matches the invalid path classification.
//
// This keeps errors.Is(err, ErrInvalidPath) stable even when callers receive
// the richer typed form.
func (e InvalidPathError) Is(target error) bool {
	return target == ErrInvalidPath
}
