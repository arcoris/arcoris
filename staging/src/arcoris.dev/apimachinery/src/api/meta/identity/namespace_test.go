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

func TestNamespace(t *testing.T) {
	namespace := Namespace("system")
	if namespace.String() != "system" {
		t.Fatalf("String() = %q", namespace.String())
	}
	if namespace.IsZero() {
		t.Fatal("non-zero Namespace IsZero() = true")
	}
	if !Namespace("").IsZero() {
		t.Fatal("zero Namespace IsZero() = false")
	}
	if !Namespace("").IsAbsent() {
		t.Fatal("zero Namespace IsAbsent() = false")
	}
	if namespace.IsAbsent() {
		t.Fatal("non-zero Namespace IsAbsent() = true")
	}
}

func TestNamespaceCanonicalText(t *testing.T) {
	text, err := Namespace("system").CanonicalText()
	requireNoError(t, err)
	if text != "system" {
		t.Fatalf("CanonicalText() = %q, want %q", text, "system")
	}

	text, err = Namespace("").CanonicalText()
	requireNoError(t, err)
	if text != "" {
		t.Fatalf("CanonicalText() = %q, want empty namespace", text)
	}

	_, err = Namespace("System").CanonicalText()
	requireErrorIs(t, err, ErrInvalidNamespace)
}
