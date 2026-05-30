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

package metagrammar

import "testing"

func TestValidateSegment(t *testing.T) {
	for _, value := range []string{"uid-1", "abc.DEF_123", "tenant:object"} {
		t.Run("valid/"+value, func(t *testing.T) {
			if err := ValidateSegment(value); err != nil {
				t.Fatalf("ValidateSegment() error = %v", err)
			}
		})
	}

	for _, value := range []string{"", "uid 1", "uid/1", "uid\n1", "uid@1"} {
		t.Run("invalid/"+value, func(t *testing.T) {
			if err := ValidateSegment(value); err == nil {
				t.Fatal("ValidateSegment() error = nil")
			}
		})
	}
}
