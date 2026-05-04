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

package health

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidTarget identifies a target that cannot be used where a concrete
	// health target is required.
	//
	// TargetUnknown is valid as a zero-value sentinel, but it is not concrete and
	// MUST NOT be used as an evaluable target, registry key, target timeout key,
	// or adapter endpoint binding.
	ErrInvalidTarget = errors.New("health: invalid target")
)

// InvalidTargetError describes a non-concrete or invalid target used at a health
// target boundary.
//
// InvalidTargetError is classified as ErrInvalidTarget. Callers should classify
// it with errors.Is and inspect Target only for diagnostics.
type InvalidTargetError struct {
	Target Target
}

// Error returns the invalid target message.
func (e InvalidTargetError) Error() string {
	return fmt.Sprintf("%v: %s", ErrInvalidTarget, e.Target.String())
}

// Is reports whether target matches the invalid target classification.
func (e InvalidTargetError) Is(target error) bool {
	return target == ErrInvalidTarget
}
