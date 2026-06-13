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

package objectstore

import (
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

func TestListItemCloneDetachesState(t *testing.T) {
	item := ListItem{Key: validKey(), State: validCommittedState()}

	cloned := item.Clone()
	cloned.State.Object.Desired = value.StringValue("mutated")

	got, ok := item.State.Object.Desired.AsString()
	if !ok || got != "desired" {
		t.Fatalf("original desired = %q, %v; want desired, true", got, ok)
	}
	if !cloned.Key.Equal(item.Key) {
		t.Fatalf("clone key = %s; want %s", cloned.Key, item.Key)
	}
}
