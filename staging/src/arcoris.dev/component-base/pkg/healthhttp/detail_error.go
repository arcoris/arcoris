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
	// ErrInvalidDetailLevel identifies an unsupported health HTTP response detail
	// level.
	//
	// The sentinel classifies only adapter exposure configuration errors.
	ErrInvalidDetailLevel = errors.New("healthhttp: invalid detail level")
)

// InvalidDetailLevelError describes an unsupported health HTTP response detail
// level.
//
// Callers should use errors.Is for classification and inspect Level only for
// configuration diagnostics.
type InvalidDetailLevelError struct {
	Level DetailLevel
}

// Error returns the invalid detail level message.
//
// The string is diagnostic-only and includes the level's stable textual form.
func (e InvalidDetailLevelError) Error() string {
	return fmt.Sprintf("%v: %s", ErrInvalidDetailLevel, e.Level.String())
}

// Is reports whether target matches the invalid detail level classification.
//
// This keeps the typed error compatible with errors.Is against
// ErrInvalidDetailLevel.
func (e InvalidDetailLevelError) Is(target error) bool {
	return target == ErrInvalidDetailLevel
}
