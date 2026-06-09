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

package types

import "testing"

func TestBytesViewAccessors(t *testing.T) {
	view := requireBytesView(t, Bytes().MinBytes(1).MaxBytes(4096).Descriptor())

	minBytes, ok := view.MinBytes()
	requireEqual(t, ok, true)
	requireEqual(t, minBytes, 1)

	maxBytes, ok := view.MaxBytes()
	requireEqual(t, ok, true)
	requireEqual(t, maxBytes, 4096)
}
