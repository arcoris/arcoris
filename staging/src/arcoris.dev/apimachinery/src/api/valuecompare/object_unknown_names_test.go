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
	"slices"
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func TestUnknownMemberNamesReturnsSortedUndeclaredNames(t *testing.T) {
	descriptor := types.Object(types.Field("known").String().Optional()).Type()
	objectView, _ := descriptor.Object()
	declared := objectFieldsByName(objectView.Fields())
	oldObject, _ := valueObject("known", "x", "zeta", "old").Object()
	newObject, _ := valueObject("alpha", "new").Object()

	got := unknownMemberNames(oldObject, newObject, declared)

	if want := []string{"alpha", "zeta"}; !slices.Equal(got, want) {
		t.Fatalf("unknownMemberNames() = %#v, want %#v", got, want)
	}
}
