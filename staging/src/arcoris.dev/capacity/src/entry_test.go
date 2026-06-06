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

func TestEntryIsValid(t *testing.T) {
	if !entry("worker_slots", 1).IsValid() {
		t.Fatal("valid entry was invalid")
	}
	if (capacity.Entry{Resource: "bad-name", Amount: 1}).IsValid() {
		t.Fatal("entry with invalid resource was valid")
	}
	if (capacity.Entry{Resource: capacity.MustResource("worker_slots")}).IsValid() {
		t.Fatal("entry with zero amount was valid")
	}
}
