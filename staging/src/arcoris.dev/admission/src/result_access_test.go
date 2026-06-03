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

package admission

import "testing"

func TestResultAccessors(t *testing.T) {
	t.Parallel()

	result := Granted(
		ReasonAdmitted,
		"lease",
		"snapshot",
	)

	if got := result.Decision(); got != Grant(ReasonAdmitted) {
		t.Fatalf("Decision = %+v, want granted admitted decision", got)
	}
	if got, ok := result.Grant(); !ok || got != "lease" {
		t.Fatalf("Grant = (%q, %v), want (lease, true)", got, ok)
	}
	if got, ok := result.Metadata(); !ok || got != "snapshot" {
		t.Fatalf("Metadata = (%q, %v), want (snapshot, true)", got, ok)
	}

	denied := DeniedFor[string](Reason("capacity_exhausted"), "snapshot")
	if got, ok := denied.Grant(); ok || got != "" {
		t.Fatalf("denied Grant = (%q, %v), want zero value and false", got, ok)
	}
}
