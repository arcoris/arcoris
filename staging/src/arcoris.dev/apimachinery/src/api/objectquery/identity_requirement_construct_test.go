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

package objectquery

import "testing"

func TestNamespaceEqualsAllowsAbsentNamespace(t *testing.T) {
	req, err := NamespaceEquals("")
	requireNoError(t, err)

	if req.IsZero() {
		t.Fatal("explicit absent namespace requirement is zero")
	}
	if req.namespace != "" {
		t.Fatalf("namespace = %q; want empty", req.namespace)
	}
}

func TestNameEqualsConstructsRequirement(t *testing.T) {
	req, err := NameEquals("worker")
	requireNoError(t, err)

	if req.IsZero() {
		t.Fatal("name requirement is zero")
	}
	if req.name != "worker" {
		t.Fatalf("name = %q; want worker", req.name)
	}
}
