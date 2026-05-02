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
	// ErrInvalidTarget identifies a target that cannot own registered checks.
	//
	// Registration requires a concrete health target. TargetUnknown is valid as a
	// zero-value sentinel, but it is not concrete and MUST NOT be used as an
	// evaluable registry key.
	ErrInvalidTarget = errors.New("health: invalid target")

	// ErrNilChecker identifies a nil checker passed to Registry.Register.
	//
	// A nil checker cannot provide a stable name or produce a Result. Registry
	// rejects nil values at registration time so evaluator code never has to
	// recover from nil entries.
	ErrNilChecker = errors.New("health: nil checker")

	// ErrDuplicateCheck identifies a repeated check name within one target.
	//
	// Check names are scoped by target. The same name MAY be registered for
	// different targets when each target has its own health semantics.
	ErrDuplicateCheck = errors.New("health: duplicate check")
)

// InvalidTargetError describes a non-concrete or invalid target used at a
// registry boundary.
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

// NilCheckerError describes a nil checker rejected from a registration batch.
//
// NilCheckerError is classified as ErrNilChecker. Target and Index identify the
// failed batch item without forcing callers to parse joined error strings.
type NilCheckerError struct {
	Target Target
	Index  int
}

// Error returns the nil checker message.
func (e NilCheckerError) Error() string {
	return fmt.Sprintf("%v: target=%s index=%d", ErrNilChecker, e.Target.String(), e.Index)
}

// Is reports whether target matches the nil checker classification.
func (e NilCheckerError) Is(target error) bool {
	return target == ErrNilChecker
}

// InvalidCheckNameError describes an invalid checker name in a registration
// batch.
//
// InvalidCheckNameError unwraps ErrEmptyCheckName or ErrInvalidCheckName so
// callers can classify the precise name failure through errors.Is even when the
// error is part of an errors.Join tree.
type InvalidCheckNameError struct {
	Target Target
	Index  int
	Name   string
	Err    error
}

// Error returns the invalid check name message.
func (e InvalidCheckNameError) Error() string {
	return fmt.Sprintf(
		"%v: target=%s index=%d name=%q",
		e.Err,
		e.Target.String(),
		e.Index,
		e.Name,
	)
}

// Unwrap returns the stable check-name validation error.
func (e InvalidCheckNameError) Unwrap() error {
	return e.Err
}

// DuplicateCheckError describes a repeated check name within one target.
//
// DuplicateCheckError is classified as ErrDuplicateCheck. Index identifies the
// duplicate item in the rejected batch. PreviousIndex is set for duplicates
// found inside the same batch; existing-registry conflicts leave it negative.
type DuplicateCheckError struct {
	Target        Target
	Name          string
	Index         int
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

// Is reports whether target matches the duplicate check classification.
func (e DuplicateCheckError) Is(target error) bool {
	return target == ErrDuplicateCheck
}
