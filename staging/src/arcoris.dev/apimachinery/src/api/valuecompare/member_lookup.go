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

import "arcoris.dev/apimachinery/api/value"

// membersByName indexes concrete record members by their string name.
func membersByName(view value.RecordView) map[string]value.Value {
	if view.IsEmpty() {
		return nil
	}

	out := make(map[string]value.Value, view.Len())
	view.ForEach(func(_ int, member value.RecordMember) bool {
		out[member.Name.String()] = member.Value
		return true
	})

	return out
}
