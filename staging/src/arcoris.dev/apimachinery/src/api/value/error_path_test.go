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

import "testing"

func TestDiagnosticPathHelpers(t *testing.T) {
	requireEqual(t, pathFloat, "float")
	requireEqual(t, pathDecimal, "decimal")
	requireEqual(t, pathDate, "date")
	requireEqual(t, pathTimeOfDay, "timeOfDay")

	requireEqual(t, objectFieldNamePath(1), "object.fields[1].name")
	requireEqual(t, objectFieldValuePath(1), "object.fields[1].value")
	requireEqual(t, listItemPath(1), "list.items[1]")
	requireEqual(t, mapEntryKeyPath(1), "map.entries[1].key")
	requireEqual(t, mapEntryValuePath(1), "map.entries[1].value")
}
