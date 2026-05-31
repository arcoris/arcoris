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

import "fmt"

const (
	// pathFloat identifies float scalar construction failures.
	pathFloat = "float"
	// pathDecimal identifies decimal scalar construction failures.
	pathDecimal = "decimal"
	// pathDate identifies date scalar construction failures.
	pathDate = "date"
	// pathTimeOfDay identifies time-of-day scalar construction failures.
	pathTimeOfDay = "timeOfDay"
)

// objectMemberNamePath returns the diagnostic path for an object member name.
//
// Paths use caller order indexes because object payloads preserve member order
// and do not store a lookup index.
func objectMemberNamePath(index int) string {
	return fmt.Sprintf("object.members[%d].name", index)
}

// objectMemberValuePath returns the diagnostic path for an object member value.
//
// The path points to the nested Value slot, not to descriptor metadata.
func objectMemberValuePath(index int) string {
	return fmt.Sprintf("object.members[%d].value", index)
}

// listItemPath returns the diagnostic path for a list item.
//
// List paths use ordered indexes because list payload order is semantically
// preserved.
func listItemPath(index int) string {
	return fmt.Sprintf("list.items[%d]", index)
}
