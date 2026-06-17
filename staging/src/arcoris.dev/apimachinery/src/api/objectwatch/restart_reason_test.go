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

package objectwatch

import "testing"

func TestRestartReasonValidity(t *testing.T) {
	tests := []struct {
		reason RestartReason
		valid  bool
	}{
		{reason: 0},
		{reason: RestartHistoryUnavailable, valid: true},
		{reason: RestartContinuityLost, valid: true},
		{reason: RestartSourceReset, valid: true},
		{reason: RestartReason(99)},
	}

	for _, tt := range tests {
		if tt.reason.IsValid() != tt.valid {
			t.Fatalf("IsValid(%d) = %v; want %v", tt.reason, tt.reason.IsValid(), tt.valid)
		}
	}
}
