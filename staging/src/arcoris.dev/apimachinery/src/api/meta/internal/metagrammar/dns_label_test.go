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

func TestValidateDNSLabel(t *testing.T) {
	for _, value := range []string{"a", "z", "0", "workers", "worker-1"} {
		t.Run("valid/"+value, func(t *testing.T) {
			if err := ValidateDNSLabel(value); err != nil {
				t.Fatalf("ValidateDNSLabel() error = %v", err)
			}
		})
	}

	for _, value := range []string{
		"",
		"Workers",
		"worker_1",
		"worker.main",
		"worker/main",
		"-worker",
		"worker-",
		"worker 1",
		"воркер",
	} {
		t.Run("invalid/"+value, func(t *testing.T) {
			if err := ValidateDNSLabel(value); err == nil {
				t.Fatal("ValidateDNSLabel() error = nil")
			}
		})
	}
}
