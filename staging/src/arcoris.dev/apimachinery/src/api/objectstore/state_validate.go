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

package objectstore

import (
	"errors"

	"arcoris.dev/apimachinery/api/objectownership"
)

// IsValid reports whether s is structurally usable as committed store state.
func (s State) IsValid() bool {
	return ValidateCommittedState(s) == nil
}

// ValidateCommittedState checks state that already has a store-assigned revision.
func ValidateCommittedState(state State) error {
	return validateState(state, true)
}

// ValidateInputState checks state supplied to Create or Update before commit.
//
// Input state must have a zero revision because concrete stores assign commit
// revisions. This prevents callers from forging concurrency tokens.
func ValidateInputState(state State) error {
	return validateState(state, false)
}

// validateState applies storage-level object and ownership checks.
func validateState(state State, committed bool) error {
	if committed && !state.Revision.IsValid() {
		return errorFor(ErrorReasonInvalidRevision, Key{}, state.Revision, 0, ErrInvalidRevision)
	}
	if !committed && !state.Revision.IsZero() {
		return errorFor(ErrorReasonInvalidRevision, Key{}, state.Revision, 0, ErrInvalidRevision)
	}

	if err := state.Object.ValidateMeta(); err != nil {
		return errorFor(ErrorReasonInvalidState, Key{}, 0, 0, errors.Join(ErrInvalidState, err))
	}
	if state.Object.Desired.IsZero() {
		return errorFor(ErrorReasonInvalidState, Key{}, 0, 0, ErrInvalidState)
	}
	if state.Object.Observed != nil && state.Object.Observed.IsZero() {
		return errorFor(ErrorReasonInvalidState, Key{}, 0, 0, ErrInvalidState)
	}
	if err := objectownership.Validate(state.Ownership); err != nil {
		return errorFor(ErrorReasonInvalidState, Key{}, 0, 0, errors.Join(ErrInvalidState, err))
	}

	return nil
}
