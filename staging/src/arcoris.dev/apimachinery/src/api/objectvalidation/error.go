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

package objectvalidation

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

// Object contract validation sentinels classify broad failure categories.
var (
	// ErrInvalidPlan classifies malformed validation plans.
	ErrInvalidPlan = errors.New("invalid object validation plan")

	// ErrInvalidObject classifies object values that do not satisfy the plan.
	ErrInvalidObject = errors.New("invalid API object")

	// ErrInvalidMetadata classifies failures returned by api/object metadata validation.
	ErrInvalidMetadata = errors.New("invalid object metadata")

	// ErrResourceMismatch classifies GVK family mismatches against the resource definition.
	ErrResourceMismatch = errors.New("object does not match resource definition")

	// ErrVersionNotDefined classifies object versions absent from the resource contract.
	ErrVersionNotDefined = errors.New("object API version is not defined by resource")

	// ErrInvalidScope classifies minimal metadata namespace/resource scope conflicts.
	ErrInvalidScope = errors.New("object metadata does not match resource scope")

	// ErrMissingValidator classifies missing desired or observed surface validators.
	ErrMissingValidator = errors.New("missing surface validator")

	// ErrInvalidDesired classifies desired surface validator failures.
	ErrInvalidDesired = errors.New("invalid desired surface")

	// ErrInvalidObserved classifies observed surface validator failures.
	ErrInvalidObserved = errors.New("invalid observed surface")

	// ErrObservedNotAllowed classifies observed values where the resource version has no observed descriptor.
	ErrObservedNotAllowed = errors.New("observed surface is not allowed")
)

// Error is the structured diagnostic returned by object contract validation.
type Error struct {
	// Record stores the shared path, sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable object validation diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("objectvalidation")
}

// Unwrap preserves broad classification and the nested validation cause.
//
// Concrete object failures also unwrap to ErrInvalidObject. Missing validators
// also unwrap to ErrInvalidPlan because they describe an incomplete plan rather
// than a malformed object value.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	base := e.Err
	if isObjectFailure(e.Err) {
		base = errors.Join(ErrInvalidObject, e.Err)
	}
	if e.Err == ErrMissingValidator {
		base = errors.Join(ErrInvalidPlan, e.Err)
	}

	return diagnostic.JoinRecord[ErrorReason](base, e.Cause).Unwrap()
}

// isObjectFailure reports whether err means the object fails the contract.
func isObjectFailure(err error) bool {
	switch err {
	case ErrInvalidMetadata,
		ErrResourceMismatch,
		ErrVersionNotDefined,
		ErrInvalidScope,
		ErrInvalidDesired,
		ErrInvalidObserved,
		ErrObservedNotAllowed:
		return true
	default:
		return false
	}
}
