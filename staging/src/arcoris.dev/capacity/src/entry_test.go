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

import "testing"

func TestEntryValidity(t *testing.T) {
	t.Parallel()

	if !entry("worker_slots", 1).IsValid() {
		t.Fatal("positive resource entry is invalid")
	}
	if entry("worker_slots", 0).IsValid() {
		t.Fatal("zero amount entry is valid")
	}
}
