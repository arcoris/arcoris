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

package healthgate

import (
	"errors"
	"fmt"

	"arcoris.dev/health"
)

var (
	// ErrInvalidGateResult identifies a result that cannot be stored in a Gate.
	//
	// Gate results must satisfy health.Result invariants. The gate may fill an
	// empty result name with its own name, but it does not accept structurally
	// invalid health observations.
	ErrInvalidGateResult = errors.New("healthgate: invalid gate result")

	// ErrMismatchedGateResult identifies a result whose non-empty name does not
	// match the owning Gate name.
	//
	// A Gate is itself a health.Checker and therefore owns one stable check name.
	// Stored results must either leave Name empty so the gate can fill it, or use
	// the exact gate name.
	ErrMismatchedGateResult = errors.New("healthgate: mismatched gate result")
)

// InvalidGateResultError describes a structurally invalid result rejected by a
// Gate.
//
// InvalidGateResultError is classified as ErrInvalidGateResult. Callers should
// use errors.Is for classification and inspect GateName or Result only for
// diagnostics.
type InvalidGateResultError struct {
	// GateName is the stable name of the gate that rejected the result.
	GateName string

	// Result is the invalid result rejected before storage.
	Result health.Result
}

// Error returns the invalid gate result message.
func (e InvalidGateResultError) Error() string {
	return fmt.Sprintf(
		"%v: gate=%q status=%s duration=%s",
		ErrInvalidGateResult,
		e.GateName,
		e.Result.Status.String(),
		e.Result.Duration,
	)
}

// Is reports whether target matches the invalid gate result classification.
func (e InvalidGateResultError) Is(target error) bool {
	return target == ErrInvalidGateResult
}

// MismatchedGateResultError describes a result whose name does not match its
// owning Gate.
//
// MismatchedGateResultError is classified as ErrMismatchedGateResult. Callers
// should use errors.Is for classification and inspect GateName or ResultName
// only for diagnostics.
type MismatchedGateResultError struct {
	// GateName is the stable name of the gate that rejected the result.
	GateName string

	// ResultName is the non-empty result name that did not match GateName.
	ResultName string
}

// Error returns the mismatched gate result message.
func (e MismatchedGateResultError) Error() string {
	return fmt.Sprintf(
		"%v: gate=%q result=%q",
		ErrMismatchedGateResult,
		e.GateName,
		e.ResultName,
	)
}

// Is reports whether target matches the mismatched gate result classification.
func (e MismatchedGateResultError) Is(target error) bool {
	return target == ErrMismatchedGateResult
}
