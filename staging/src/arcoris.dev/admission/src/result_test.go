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

func TestZeroResultIsInvalid(t *testing.T) {
	t.Parallel()

	var result Result[NoGrant, NoMetadata]
	if result.IsValid() {
		t.Fatal("zero result should be invalid")
	}
}

func TestResultStateHelpers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		result     Result[string, string]
		admitted   bool
		denied     bool
		queued     bool
		deferred   bool
		sideEffect bool
		grant      bool
		metadata   bool
	}{
		{
			name: "queued with handle",
			result: Queued(
				ReasonQueued,
				"ticket",
				"snapshot",
			),
			queued:     true,
			sideEffect: true,
			grant:      true,
			metadata:   true,
		},
		{
			name:     "deferred with metadata",
			result:   DeferredFor[string](ReasonDeferred, "snapshot"),
			deferred: true,
			metadata: true,
		},
		{
			name: "granted no metadata",
			result: resultWith(
				Grant(ReasonAdmitted),
				someString("lease"),
				noneString(),
			),
			admitted:   true,
			sideEffect: true,
			grant:      true,
		},
		{
			name: "denied no metadata",
			result: resultWith(
				Deny(Reason("capacity_exhausted")),
				noneString(),
				noneString(),
			),
			denied: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.result.IsAdmitted(); got != tt.admitted {
				t.Fatalf("IsAdmitted = %v, want %v", got, tt.admitted)
			}
			if got := tt.result.IsDenied(); got != tt.denied {
				t.Fatalf("IsDenied = %v, want %v", got, tt.denied)
			}
			if got := tt.result.IsQueued(); got != tt.queued {
				t.Fatalf("IsQueued = %v, want %v", got, tt.queued)
			}
			if got := tt.result.IsDeferred(); got != tt.deferred {
				t.Fatalf("IsDeferred = %v, want %v", got, tt.deferred)
			}
			if got := tt.result.HasSideEffect(); got != tt.sideEffect {
				t.Fatalf("HasSideEffect = %v, want %v", got, tt.sideEffect)
			}
			if got := tt.result.HasGrant(); got != tt.grant {
				t.Fatalf("HasGrant = %v, want %v", got, tt.grant)
			}
			if got := tt.result.HasMetadata(); got != tt.metadata {
				t.Fatalf("HasMetadata = %v, want %v", got, tt.metadata)
			}
		})
	}
}
