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

import "errors"

// ValidateChange checks the invariant set for every supported transition kind.
func ValidateChange(change Change) error {
	if !change.Kind.IsValid() {
		return changeError(
			ErrorReasonInvalidChangeKind,
			change.Key,
			0,
			0,
			nil,
		)
	}
	if err := ValidateKey(change.Key); err != nil {
		return changeError(
			ErrorReasonInvalidChange,
			change.Key,
			0,
			0,
			err,
		)
	}
	if !change.Revision.IsValid() {
		return changeError(
			ErrorReasonInvalidChange,
			change.Key,
			change.Revision,
			0,
			ErrInvalidRevision,
		)
	}

	switch change.Kind {
	case ChangeCreated:
		return validateCreatedChange(change)
	case ChangeUpdated:
		return validateUpdatedChange(change)
	case ChangeDeleted:
		return validateDeletedChange(change)
	default:
		return changeError(ErrorReasonInvalidChangeKind, change.Key, 0, 0, nil)
	}
}

// validateCreatedChange checks missing/deleted -> live transition shape.
func validateCreatedChange(change Change) error {
	if !isZeroState(change.Before) {
		return changeError(
			ErrorReasonInvalidChange,
			change.Key,
			0,
			change.Before.Revision,
			nil,
		)
	}
	if err := ValidateCommittedState(change.After); err != nil {
		return changeError(
			ErrorReasonInvalidChange,
			change.Key,
			change.Revision,
			change.After.Revision,
			err,
		)
	}
	if change.After.Revision != change.Revision {
		return changeError(
			ErrorReasonInvalidChange,
			change.Key,
			change.Revision,
			change.After.Revision,
			ErrInvalidRevision,
		)
	}

	return nil
}

// validateUpdatedChange checks live -> live transition shape.
func validateUpdatedChange(change Change) error {
	if err := ValidateCommittedState(change.Before); err != nil {
		return changeError(
			ErrorReasonInvalidChange,
			change.Key,
			change.Revision,
			change.Before.Revision,
			err,
		)
	}
	if err := ValidateCommittedState(change.After); err != nil {
		return changeError(
			ErrorReasonInvalidChange,
			change.Key,
			change.Revision,
			change.After.Revision,
			err,
		)
	}
	if change.After.Revision != change.Revision {
		return changeError(
			ErrorReasonInvalidChange,
			change.Key,
			change.Revision,
			change.After.Revision,
			ErrInvalidRevision,
		)
	}
	if !change.Before.Revision.Before(change.After.Revision) {
		return changeError(
			ErrorReasonInvalidChange,
			change.Key,
			change.Before.Revision,
			change.After.Revision,
			ErrInvalidRevision,
		)
	}

	return nil
}

// validateDeletedChange checks live -> tombstone transition shape.
func validateDeletedChange(change Change) error {
	if err := ValidateCommittedState(change.Before); err != nil {
		return changeError(
			ErrorReasonInvalidChange,
			change.Key,
			change.Revision,
			change.Before.Revision,
			err,
		)
	}
	if !isZeroState(change.After) {
		return changeError(
			ErrorReasonInvalidChange,
			change.Key,
			0,
			change.After.Revision,
			nil,
		)
	}
	if !change.Before.Revision.Before(change.Revision) {
		return changeError(
			ErrorReasonInvalidChange,
			change.Key,
			change.Before.Revision,
			change.Revision,
			ErrInvalidRevision,
		)
	}

	return nil
}

// changeError wraps lower causes with ErrInvalidChange as the broad category.
func changeError(reason ErrorReason, key Key, expected Revision, actual Revision, cause error) error {
	err := ErrInvalidChange
	if cause != nil {
		err = errors.Join(ErrInvalidChange, cause)
	}

	return errorFor(reason, key, expected, actual, err)
}

// isZeroState reports whether state is the intentionally absent transition side.
func isZeroState(state State) bool {
	return state.Revision.IsZero() &&
		state.Object.TypeMeta.IsZero() &&
		state.Object.ObjectMeta.IsZero() &&
		state.Object.Desired.IsZero() &&
		state.Object.Observed == nil &&
		state.Ownership.IsEmpty()
}
