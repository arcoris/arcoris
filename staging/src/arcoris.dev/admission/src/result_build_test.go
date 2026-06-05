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

func TestResultWithBuildsExactShape(t *testing.T) {
	t.Parallel()

	decision := GrantDecision(ReasonAdmitted)
	result := resultWith(decision, "grant", true, "metadata", true)

	if got := result.Decision(); got != decision {
		t.Fatalf("Decision() = %+v, want %+v", got, decision)
	}
	if got, ok := result.Grant(); !ok || got != "grant" {
		t.Fatalf("Grant() = (%q, %t), want grant,true", got, ok)
	}
	if got, ok := result.Metadata(); !ok || got != "metadata" {
		t.Fatalf("Metadata() = (%q, %t), want metadata,true", got, ok)
	}
}

func TestResultWithCanBuildIntentionallyInvalidShapes(t *testing.T) {
	t.Parallel()

	result := resultWith(GrantDecision(ReasonAdmitted), "", false, NoMetadata{}, false)
	if result.IsValid() {
		t.Fatal("owned result without grant is valid")
	}
}
