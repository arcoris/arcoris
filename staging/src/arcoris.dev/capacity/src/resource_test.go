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
	"errors"
	"testing"

	"arcoris.dev/capacity"
)

func TestResourceValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value string
		valid bool
	}{
		{value: "worker_slots", valid: true},
		{value: "resilience.bulkhead.slots", valid: true},
		{value: "scheduler.node.cpu_units", valid: true},
		{value: "", valid: false},
		{value: "Worker_slots", valid: false},
		{value: "worker-slots", valid: false},
		{value: "worker slots", valid: false},
		{value: ".worker_slots", valid: false},
		{value: "worker_slots.", valid: false},
		{value: "worker..slots", valid: false},
		{value: "_worker_slots", valid: false},
		{value: "worker_slots_", valid: false},
		{value: "worker__slots", valid: false},
		{value: "worker.1slots", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			t.Parallel()
			if got := capacity.Resource(tt.value).IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestMustResourcePanicsWithStructuredError(t *testing.T) {
	t.Parallel()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("MustResource did not panic")
		}
		err, ok := recovered.(error)
		if !ok {
			t.Fatalf("panic = %#v, want error", recovered)
		}
		if !errors.Is(err, capacity.ErrInvalidResource) {
			t.Fatalf("panic error = %v, want ErrInvalidResource", err)
		}
		var capacityErr *capacity.Error
		if !errors.As(err, &capacityErr) {
			t.Fatalf("panic error = %T, want *capacity.Error", err)
		}
		if capacityErr.Path != "resource" {
			t.Fatalf("Path = %q, want resource", capacityErr.Path)
		}
	}()

	_ = capacity.MustResource("bad-name")
}
