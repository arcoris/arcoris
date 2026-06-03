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

package valuecompare

import "testing"

func TestListMapOperandKeepsPresence(t *testing.T) {
	entry := listMapEntry{item: conditionValue("Ready", "True")}

	got := listMapOperand(entry, true)
	equal, err := newComparer(Options{}).equalOpaqueValue(rootField("conditions").Index(0), got.value, entry.item)
	requireNoError(t, err)
	if !got.present || !equal {
		t.Fatalf("listMapOperand(present) = %#v", got)
	}

	got = listMapOperand(entry, false)
	equal, err = newComparer(Options{}).equalOpaqueValue(rootField("conditions").Index(0), got.value, entry.item)
	requireNoError(t, err)
	if got.present || !equal {
		t.Fatalf("listMapOperand(absent) = %#v", got)
	}
}
