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
type ConflictError struct {
	// Conflicts contains the deterministic conflict details.
	Conflicts ConflictSet
}

// Error returns deterministic conflict text.
func (e *ConflictError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Conflicts.Error()
}

// Is reports ErrConflict identity for errors.Is.
func (e *ConflictError) Is(target error) bool {
	return target == ErrConflict
}
