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

func TestMapKey(t *testing.T) {
	key, err := NewMapKey("app")

	requireNoError(t, err)
	requireEqual(t, key.String(), "app")
	requireEqual(t, key.IsZero(), false)
}

func TestMapKeyRejectsEmptyKey(t *testing.T) {
	_, err := NewMapKey("")

	requireErrorIs(t, err, ErrEmptyMapKey)
}

func TestMustMapKeyPanicsOnEmptyKey(t *testing.T) {
	requirePanic(t, func() {
		MustMapKey("")
	})
}
