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
