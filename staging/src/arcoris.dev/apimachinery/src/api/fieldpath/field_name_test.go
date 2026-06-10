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

func TestFieldName(t *testing.T) {
	name, err := NewFieldName("spec")

	requireNoError(t, err)
	requireEqual(t, name.String(), "spec")
	requireEqual(t, name.IsZero(), false)
}

func TestFieldNameRejectsEmptyName(t *testing.T) {
	_, err := NewFieldName("")

	requireErrorIs(t, err, ErrEmptyFieldName)
}

func TestMustFieldNamePanicsOnEmptyName(t *testing.T) {
	requirePanic(t, func() {
		MustFieldName("")
	})
}
