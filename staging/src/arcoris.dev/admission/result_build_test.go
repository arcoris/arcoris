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

func TestResultWithPreservesDecisionGrantAndMetadata(t *testing.T) {
	t.Parallel()

	result := resultWith(
		Grant(ReasonAdmitted),
		someString("lease"),
		someString("snapshot"),
	)

	if result.Decision() != Grant(ReasonAdmitted) {
		t.Fatalf("decision = %+v, want granted admitted decision", result.Decision())
	}
	if grant, ok := result.Grant(); !ok || grant != "lease" {
		t.Fatalf("grant = (%q, %v), want (lease, true)", grant, ok)
	}
	if metadata, ok := result.Metadata(); !ok || metadata != "snapshot" {
		t.Fatalf("metadata = (%q, %v), want (snapshot, true)", metadata, ok)
	}
}
