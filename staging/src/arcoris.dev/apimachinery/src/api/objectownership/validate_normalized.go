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

package objectownership

import "reflect"

// ValidateNormalized checks that state is valid and already canonical.
//
// The check compares state to Normalize(state). It does not stamp or inspect a
// document version because State is the only ownership model. Stores use this to
// assert that committed state is already in canonical form.
func ValidateNormalized(state State) error {
	if err := Validate(state); err != nil {
		return err
	}
	if !reflect.DeepEqual(state, Normalize(state)) {
		return errorAt(
			pathState,
			ErrNotNormalized,
			ErrorReasonNotNormalized,
			"object ownership state is not normalized",
		)
	}

	return nil
}
