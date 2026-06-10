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

package valuemerge

import (
	"testing"

	"arcoris.dev/apimachinery/api/internal/valuepresence"
)

func TestListItemsReturnsDetachedItems(t *testing.T) {
	got := listItems(valuepresence.Present(list(str("a"))))

	if len(got) != 1 {
		t.Fatalf("items length = %d; want 1", len(got))
	}
	text, _ := got[0].AsString()
	if text != "a" {
		t.Fatalf("item = %q; want a", text)
	}
}

func TestItemAtAbsentOutOfRange(t *testing.T) {
	if itemAt(nil, 0).Present() {
		t.Fatalf("itemAt(nil, 0) is present")
	}
}

func TestAppendItemSkipsAbsent(t *testing.T) {
	got := appendItem(nil, valuepresence.Absent())

	if len(got) != 0 {
		t.Fatalf("items length = %d; want 0", len(got))
	}
}
