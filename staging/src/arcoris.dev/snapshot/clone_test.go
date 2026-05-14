/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package snapshot

import "testing"

func TestIdentity(t *testing.T) {
	if got, want := Identity("value"), "value"; got != want {
		t.Fatalf("Identity() = %q, want %q", got, want)
	}
}

func TestRequireClonePanicsOnNil(t *testing.T) {
	requirePanicWith(t, "snapshot: nil clone function", func() {
		_ = requireClone[string](nil)
	})
}

func TestRequireCloneReturnsClone(t *testing.T) {
	clone := func(v string) string { return v }
	if requireClone(clone)("value") != "value" {
		t.Fatal("requireClone did not return supplied clone")
	}
}
