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

package fieldownership

// ConflictError wraps a non-empty conflict set for errors.Is support.
//
// Prefer NewConflictError when converting conflicts into an error. It returns
// nil for empty conflict sets and stores non-empty conflicts in deterministic
// order.
type ConflictError struct {
	conflicts ConflictSet
}

// NewConflictError returns nil for empty conflicts or a deterministic conflict error.
func NewConflictError(conflicts ConflictSet) error {
	if conflicts.IsEmpty() {
		return nil
	}

	return &ConflictError{
		conflicts: newConflictSetUnchecked(conflicts.Conflicts()...),
	}
}

// Error returns deterministic conflict text.
func (e *ConflictError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.conflicts.Error()
}

// Is reports ErrConflict identity for errors.Is.
func (e *ConflictError) Is(target error) bool {
	return target == ErrConflict
}

// Conflicts returns detached deterministic conflict details.
func (e *ConflictError) Conflicts() ConflictSet {
	if e == nil {
		return ConflictSet{}
	}

	return newConflictSetUnchecked(e.conflicts.Conflicts()...)
}
