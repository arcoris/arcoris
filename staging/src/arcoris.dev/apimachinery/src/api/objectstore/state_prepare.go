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

// PrepareInputState validates, detaches, and canonicalizes Create/Update state.
//
// The returned state still has zero revision. Concrete stores assign a
// committed revision only after their concurrency preconditions pass.
func PrepareInputState(state State) (State, error) {
	if err := ValidateInputState(state); err != nil {
		return State{}, err
	}

	prepared := state.Clone()
	normalized, err := objectownership.Normalize(prepared.Ownership)
	if err != nil {
		return State{}, errorFor(ReasonInvalidState, Key{}, 0, 0, errors.Join(ErrInvalidState, err))
	}
	prepared.Ownership = normalized

	return prepared, nil
}
