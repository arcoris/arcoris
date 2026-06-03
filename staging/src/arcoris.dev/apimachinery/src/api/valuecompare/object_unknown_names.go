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

package valuecompare

import (
	"slices"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// unknownMemberNames returns sorted undeclared names present in either object.
func unknownMemberNames(
	oldObject value.ObjectView,
	newObject value.ObjectView,
	declared map[string]types.FieldDescriptor,
) []string {
	seen := make(map[string]bool, oldObject.Len()+newObject.Len())
	addUnknownNames(seen, oldObject, declared)
	addUnknownNames(seen, newObject, declared)

	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	slices.Sort(names)

	return names
}

// addUnknownNames adds undeclared object member names to seen.
func addUnknownNames(
	seen map[string]bool,
	object value.ObjectView,
	declared map[string]types.FieldDescriptor,
) {
	for _, member := range object.Members() {
		if _, ok := declared[member.Name]; !ok {
			seen[member.Name] = true
		}
	}
}
