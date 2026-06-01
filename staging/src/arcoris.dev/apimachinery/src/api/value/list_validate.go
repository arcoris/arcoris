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

package value

// validateListItem rejects uninitialized list payload values.
//
// Null is accepted because it is explicit payload data. Only KindInvalid is
// rejected as a missing construction signal.
func validateListItem(index int, item Value) error {
	if !item.IsZero() {
		return nil
	}

	return newError(
		listItemPath(index),
		ErrInvalidList,
		ErrorReasonInvalidValue,
		"list item has an invalid zero value",
	)
}
