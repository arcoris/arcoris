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

func TestName(t *testing.T) {
	name := Name("worker")
	if name.String() != "worker" {
		t.Fatalf("String() = %q", name.String())
	}
	if name.IsZero() {
		t.Fatal("non-zero Name IsZero() = true")
	}
	if !Name("").IsZero() {
		t.Fatal("zero Name IsZero() = false")
	}
	if !Name("").IsAbsent() {
		t.Fatal("zero Name IsAbsent() = false")
	}
	if name.IsAbsent() {
		t.Fatal("non-zero Name IsAbsent() = true")
	}
}

func TestNameCanonicalText(t *testing.T) {
	text, err := Name("worker").CanonicalText()
	requireNoError(t, err)
	if text != "worker" {
		t.Fatalf("CanonicalText() = %q, want %q", text, "worker")
	}

	_, err = Name("Worker").CanonicalText()
	requireErrorIs(t, err, ErrInvalidName)
}
