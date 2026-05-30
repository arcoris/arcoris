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

func TestValidateMapValue(t *testing.T) {
	strict := MapValueOptions{AllowEmpty: true, MaxLength: 63, Strict: true}
	valid := []string{"", "worker", "worker-1", "worker_1", "worker.1", "1"}

	for _, value := range valid {
		t.Run("valid/"+value, func(t *testing.T) {
			if err := ValidateMapValue(value, strict); err != nil {
				t.Fatalf("ValidateMapValue(%q) error = %v", value, err)
			}
		})
	}

	invalid := []string{
		"worker 1",
		"-worker",
		"worker-",
		"_worker",
		"worker_",
		".worker",
		"worker.",
		".",
		"---",
		"___",
	}

	for _, value := range invalid {
		t.Run("invalid/"+value, func(t *testing.T) {
			if err := ValidateMapValue(value, strict); err == nil {
				t.Fatalf("ValidateMapValue(%q) error = nil", value)
			}
		})
	}

	loose := MapValueOptions{AllowEmpty: true}
	if err := ValidateMapValue("worker 1", loose); err != nil {
		t.Fatalf("ValidateMapValue(loose) error = %v", err)
	}
	if err := ValidateMapValue("worker\n1", loose); err == nil {
		t.Fatal("ValidateMapValue(loose) error = nil")
	}
}
