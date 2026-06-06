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

func TestRefusalValidityAndPredicates(t *testing.T) {
	valid := []capacity.Refusal{
		capacity.RefusalNone,
		capacity.RefusalInsufficient,
		capacity.RefusalDebt,
		capacity.RefusalUnknownResource,
	}
	for _, refusal := range valid {
		if !refusal.IsValid() {
			t.Fatalf("%s was invalid", refusal)
		}
	}

	if capacity.RefusalNone.Refused() {
		t.Fatal("RefusalNone.Refused() = true")
	}
	if !capacity.RefusalDebt.Refused() {
		t.Fatal("RefusalDebt.Refused() = false")
	}
	if capacity.Refusal(99).IsValid() {
		t.Fatal("unknown refusal was valid")
	}
}
