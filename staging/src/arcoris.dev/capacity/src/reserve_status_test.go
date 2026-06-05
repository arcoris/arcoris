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

package capacity_test

import (
	"testing"

	"arcoris.dev/capacity"
)

func TestReserveStatusPredicatesAndStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status   capacity.ReserveStatus
		valid    bool
		reserved bool
		denied   bool
		text     string
	}{
		{status: 0, text: "invalid"},
		{status: capacity.ReserveStatusReserved, valid: true, reserved: true, text: "reserved"},
		{status: capacity.ReserveStatusInsufficient, valid: true, denied: true, text: "insufficient"},
		{status: capacity.ReserveStatusDebt, valid: true, denied: true, text: "debt"},
		{status: capacity.ReserveStatusUnknownResource, valid: true, denied: true, text: "unknown_resource"},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			t.Parallel()
			if got := tt.status.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v, want %v", got, tt.valid)
			}
			if got := tt.status.Reserved(); got != tt.reserved {
				t.Fatalf("Reserved() = %v, want %v", got, tt.reserved)
			}
			if got := tt.status.Denied(); got != tt.denied {
				t.Fatalf("Denied() = %v, want %v", got, tt.denied)
			}
			if got := tt.status.String(); got != tt.text {
				t.Fatalf("String() = %q, want %q", got, tt.text)
			}
		})
	}
}
