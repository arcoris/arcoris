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

package fieldpath

import "testing"

func TestElementKindString(t *testing.T) {
	requireEqual(t, ElementField.String(), "field")
	requireEqual(t, ElementKey.String(), "key")
	requireEqual(t, ElementIndex.String(), "index")
	requireEqual(t, ElementSelector.String(), "selector")
}

func TestElementKindIsValid(t *testing.T) {
	requireEqual(t, ElementInvalid.IsValid(), false)
	requireEqual(t, ElementField.IsValid(), true)
	requireEqual(t, ElementKey.IsValid(), true)
	requireEqual(t, ElementIndex.IsValid(), true)
	requireEqual(t, ElementSelector.IsValid(), true)
}
