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

func TestResourceValidation(t *testing.T) {
	tests := map[string]bool{
		"worker_slots":                 true,
		"runtime.worker_pool.workers":  true,
		"retry_budget_units2":          true,
		"":                             false,
		"Worker_slots":                 false,
		"worker-slots":                 false,
		"worker slots":                 false,
		".worker_slots":                false,
		"worker_slots.":                false,
		"worker..slots":                false,
		"_worker_slots":                false,
		"worker_slots_":                false,
		"worker__slots":                false,
		"worker_slots.9queue":          false,
		"runtime.worker_pool.request1": true,
	}

	for value, valid := range tests {
		if got := capacity.Resource(value).IsValid(); got != valid {
			t.Fatalf("Resource(%q).IsValid() = %v, want %v", value, got, valid)
		}
	}
}

func TestMustResourcePanicsOnInvalidValue(t *testing.T) {
	requirePanicIs(t, capacity.ErrInvalidResource, func() {
		_ = capacity.MustResource("bad-name")
	})
}
