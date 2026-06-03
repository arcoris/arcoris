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

package typeref

import "testing"

func TestFailureBuildsError(t *testing.T) {
	err := failure(rootPath(), FailureInvalidDescriptor, "bad descriptor")

	refError, ok := AsError(err)
	if !ok {
		t.Fatalf("failure() error = %v; want typeref.Error", err)
	}
	if refError.Path.String() != "$" {
		t.Fatalf("path = %s; want $", refError.Path)
	}
	if refError.Kind != FailureInvalidDescriptor {
		t.Fatalf("kind = %s; want %s", refError.Kind, FailureInvalidDescriptor)
	}
	if refError.Detail != "bad descriptor" {
		t.Fatalf("detail = %q", refError.Detail)
	}
}

func TestFailurefFormatsDetail(t *testing.T) {
	err := failuref(rootPath(), FailureUnresolvedRef, "reference %q was not found", "example.Name")

	refError, ok := AsError(err)
	if !ok {
		t.Fatalf("failuref() error = %v; want typeref.Error", err)
	}
	if refError.Detail != `reference "example.Name" was not found` {
		t.Fatalf("detail = %q", refError.Detail)
	}
}
