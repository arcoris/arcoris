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

package valueapply

import "testing"

func TestResultZeroValue(t *testing.T) {
	result := Result{}

	if !result.Value.IsZero() {
		t.Fatalf("Value is not zero")
	}
	if !result.Ownership.IsEmpty() {
		t.Fatalf("Ownership is not empty")
	}
	if !result.AppliedFields.IsEmpty() {
		t.Fatalf("AppliedFields is not empty")
	}
	if !result.IsZero() {
		t.Fatalf("IsZero() = false")
	}
	if result.HasValue() {
		t.Fatalf("HasValue() = true")
	}
	if result.HasOwnership() {
		t.Fatalf("HasOwnership() = true")
	}
}

func TestResultHasValueAndOwnership(t *testing.T) {
	result := Result{
		Value:     str("merged"),
		Ownership: state(entry("user", root())),
	}

	if !result.HasValue() {
		t.Fatalf("HasValue() = false")
	}
	if !result.HasOwnership() {
		t.Fatalf("HasOwnership() = false")
	}
	if result.IsZero() {
		t.Fatalf("IsZero() = true")
	}
}
