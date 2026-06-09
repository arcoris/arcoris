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

package identity

import "testing"

func TestGroupValue(t *testing.T) {
	requireString(t, Group("control.arcoris.dev").String(), "control.arcoris.dev")

	if !CoreGroup.IsZero() {
		t.Fatalf("CoreGroup.IsZero() = false, want true")
	}
	if !CoreGroup.IsCore() {
		t.Fatalf("CoreGroup.IsCore() = false, want true")
	}
	if Group("control.arcoris.dev").IsZero() {
		t.Fatalf("named group should not be zero")
	}
	if Group("control.arcoris.dev").IsCore() {
		t.Fatalf("named group should not be core")
	}
}
