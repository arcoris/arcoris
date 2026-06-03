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

import (
	"arcoris.dev/apimachinery/api/value"
	"slices"
	"testing"
)

func TestUnionSortedKeys(t *testing.T) {
	got := unionSortedKeys(
		map[string]value.Value{"b": value.StringValue("b"), "a": value.StringValue("a")},
		map[string]value.Value{"c": value.StringValue("c"), "a": value.StringValue("a")},
	)

	if want := []string{"a", "b", "c"}; !slices.Equal(got, want) {
		t.Fatalf("unionSortedKeys() = %#v, want %#v", got, want)
	}
}
