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

package health

import (
	"errors"
	"fmt"
)

// ErrDuplicateCheckName identifies duplicate check names in one check set.
//
// Check names are unique within a concrete target so reports, diagnostics, and
// adapters can address individual checks deterministically. The same name may be
// used by separate targets when each target owns different health semantics.
var ErrDuplicateCheckName = errors.New("health: duplicate check name")

// NilCheckError describes a nil checker in a CheckSet constructor input.
//
// NilCheckError is classified as ErrNilChecker. Index identifies the rejected
// checker position in the caller-supplied input.
type NilCheckError struct {
	// Index is the checker position in the constructor input.
	Index int
}

// Error returns the nil check message.
func (e NilCheckError) Error() string {
	return fmt.Sprintf("%v: index=%d", ErrNilChecker, e.Index)
}

// Is reports whether target matches ErrNilChecker.
func (e NilCheckError) Is(target error) bool {
	return target == ErrNilChecker
}

// InvalidCheckNameError describes an invalid checker name in a CheckSet
// constructor input.
//
// InvalidCheckNameError unwraps ErrEmptyCheckName or ErrInvalidCheckName. Index
// identifies the rejected checker position in the caller-supplied input.
type InvalidCheckNameError struct {
	// Index is the checker position in the constructor input.
	Index int

	// Name is the invalid checker name returned by Checker.Name.
	Name string

	// Err is the root check-name validation error.
	Err error
}

// Error returns the invalid check-name message.
func (e InvalidCheckNameError) Error() string {
	return fmt.Sprintf("%v: index=%d name=%q", e.Err, e.Index, e.Name)
}

// Unwrap returns the root check-name validation error.
func (e InvalidCheckNameError) Unwrap() error {
	return e.Err
}

// DuplicateCheckNameError describes a duplicate checker name in a CheckSet.
//
// The error is classified as ErrDuplicateCheckName. Index identifies the later
// checker and PreviousIndex identifies the earlier checker with the same name.
type DuplicateCheckNameError struct {
	// Name is the duplicated checker name.
	Name string

	// Index is the later checker position in the constructor input.
	Index int

	// PreviousIndex is the earlier checker position with the same name.
	PreviousIndex int
}

// Error returns the duplicate check-name message.
func (e DuplicateCheckNameError) Error() string {
	return fmt.Sprintf(
		"%v: name=%q index=%d previous_index=%d",
		ErrDuplicateCheckName,
		e.Name,
		e.Index,
		e.PreviousIndex,
	)
}

// Is reports whether target matches ErrDuplicateCheckName.
func (e DuplicateCheckNameError) Is(target error) bool {
	return target == ErrDuplicateCheckName
}
