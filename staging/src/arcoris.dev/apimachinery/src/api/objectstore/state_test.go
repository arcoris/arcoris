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

package objectstore

import "testing"

func TestStateCarriesObjectOwnershipAndRevision(t *testing.T) {
	state := validCommittedState()
	state.Ownership = ownershipWithEntry()

	if string(state.Object.ObjectMeta.Name) != "main" {
		t.Fatalf("object name = %q; want %q", state.Object.ObjectMeta.Name, "main")
	}
	if len(state.Ownership.Desired.Entries) != 1 {
		t.Fatalf("ownership entries = %d; want 1", len(state.Ownership.Desired.Entries))
	}
	if state.Revision != 1 {
		t.Fatalf("Revision = %v; want 1", state.Revision)
	}
}
