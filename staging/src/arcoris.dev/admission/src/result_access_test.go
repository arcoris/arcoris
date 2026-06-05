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

func TestResultAccessorsReturnDecisionGrantAndMetadata(t *testing.T) {
	t.Parallel()

	result := GrantedResult(ReasonAdmitted, "lease", "snapshot")
	if got, want := result.Decision(), GrantDecision(ReasonAdmitted); got != want {
		t.Fatalf("Decision() = %+v, want %+v", got, want)
	}
	if got, ok := result.Grant(); !ok || got != "lease" {
		t.Fatalf("Grant() = (%q, %t), want lease,true", got, ok)
	}
	if got, ok := result.Metadata(); !ok || got != "snapshot" {
		t.Fatalf("Metadata() = (%q, %t), want snapshot,true", got, ok)
	}
}

func TestResultAbsentGrantReturnsZeroFalse(t *testing.T) {
	t.Parallel()

	result := DeniedForResult[*struct{}](ReasonDenied, "snapshot")
	if got, ok := result.Grant(); ok || got != nil {
		t.Fatalf("Grant() = (%v, %t), want nil,false", got, ok)
	}
}

func TestResultAbsentMetadataReturnsZeroFalse(t *testing.T) {
	t.Parallel()

	result := GrantedNoMetadataResult(ReasonAdmitted, "lease")
	if got, ok := result.Metadata(); ok || got != (NoMetadata{}) {
		t.Fatalf("Metadata() = (%+v, %t), want zero,false", got, ok)
	}
}

func TestZeroResultIsInvalid(t *testing.T) {
	t.Parallel()

	var result Result[string, string]
	if result.IsValid() {
		t.Fatal("zero-value Result is valid, want invalid")
	}
}
