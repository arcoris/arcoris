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

// objectFieldNamePath returns the diagnostic path for an object field name.
//
// Paths use caller order indexes because object payloads preserve field order
// and do not store a lookup index.
func objectFieldNamePath(index int) string {
	return fmt.Sprintf("object.fields[%d].name", index)
}

// objectFieldValuePath returns the diagnostic path for an object field value.
//
// The path points to the nested Value slot, not to descriptor metadata.
func objectFieldValuePath(index int) string {
	return fmt.Sprintf("object.fields[%d].value", index)
}

// listItemPath returns the diagnostic path for a list item.
//
// List paths use ordered indexes because list payload order is semantically
// preserved.
func listItemPath(index int) string {
	return fmt.Sprintf("list.items[%d]", index)
}

// mapEntryKeyPath returns the diagnostic path for a map entry key.
//
// Map entries preserve caller order, so diagnostics identify the entry position
// that failed construction.
func mapEntryKeyPath(index int) string {
	return fmt.Sprintf("map.entries[%d].key", index)
}

// mapEntryValuePath returns the diagnostic path for a map entry value.
//
// The path points to the nested Value slot, not to future map value descriptors.
func mapEntryValuePath(index int) string {
	return fmt.Sprintf("map.entries[%d].value", index)
}
