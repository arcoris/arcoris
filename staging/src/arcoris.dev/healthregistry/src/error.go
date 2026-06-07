// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package healthregistry

import (
	"errors"
	"fmt"

	"arcoris.dev/health"
)

// ErrDuplicateCheck identifies a repeated check name within one target.
var ErrDuplicateCheck = errors.New("healthregistry: duplicate check")

// NilCheckerError describes a nil checker rejected from a registration batch.
//
// The error is classified as health.ErrNilChecker. Target and Index identify the
// rejected batch item.
type NilCheckerError struct {
	// Target is the registration target that received the nil checker.
	Target health.Target

	// Index is the checker position in the rejected batch.
	Index int
}

// Error returns the nil checker message.
func (e NilCheckerError) Error() string {
	return fmt.Sprintf("%v: target=%s index=%d", health.ErrNilChecker, e.Target.String(), e.Index)
}

// Is reports whether target matches health.ErrNilChecker.
func (e NilCheckerError) Is(target error) bool {
	return target == health.ErrNilChecker
}

// InvalidCheckNameError describes an invalid checker name in a registration
// batch.
//
// The error unwraps health.ErrEmptyCheckName or health.ErrInvalidCheckName.
type InvalidCheckNameError struct {
	// Target is the registration target that received the invalid checker.
	Target health.Target

	// Index is the checker position in the rejected batch.
	Index int

	// Name is the invalid checker name returned by Checker.Name.
	Name string

	// Err is the root health check-name validation error.
	Err error
}

// Error returns the invalid check-name message.
func (e InvalidCheckNameError) Error() string {
	return fmt.Sprintf(
		"%v: target=%s index=%d name=%q",
		e.Err,
		e.Target.String(),
		e.Index,
		e.Name,
	)
}

// Unwrap returns the root check-name validation error.
func (e InvalidCheckNameError) Unwrap() error {
	return e.Err
}

// DuplicateCheckError describes a duplicate check name within one target.
//
// The error is classified as ErrDuplicateCheck. Index identifies the duplicate
// item in the rejected batch. PreviousIndex identifies the first check with the
// same name under the target, whether it came from the same batch or from an
// earlier successful registration.
type DuplicateCheckError struct {
	// Target is the target where the duplicate name was detected.
	Target health.Target

	// Name is the duplicated checker name.
	Name string

	// Index is the later checker position in the rejected batch.
	Index int

	// PreviousIndex is the earlier checker position under the same target.
	PreviousIndex int
}

// Error returns the duplicate check message.
func (e DuplicateCheckError) Error() string {
	if e.PreviousIndex >= 0 {
		return fmt.Sprintf(
			"%v: target=%s name=%q index=%d previous_index=%d",
			ErrDuplicateCheck,
			e.Target.String(),
			e.Name,
			e.Index,
			e.PreviousIndex,
		)
	}

	return fmt.Sprintf(
		"%v: target=%s name=%q index=%d",
		ErrDuplicateCheck,
		e.Target.String(),
		e.Name,
		e.Index,
	)
}

// Is reports whether target matches ErrDuplicateCheck.
func (e DuplicateCheckError) Is(target error) bool {
	return target == ErrDuplicateCheck
}
