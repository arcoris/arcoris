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

package meta

import "testing"

func TestPageMetaClone(t *testing.T) {
	count := uint64(3)
	meta := PageMeta{RemainingItemCount: &count}

	cloned := meta.Clone()
	*cloned.RemainingItemCount = 9

	if *meta.RemainingItemCount != 3 {
		t.Fatal("RemainingItemCount pointer was not detached")
	}
	*meta.RemainingItemCount = 4
	if *cloned.RemainingItemCount != 9 {
		t.Fatal("original RemainingItemCount mutation changed clone")
	}
}

func TestPageMetaCloneNilRemainingItemCount(t *testing.T) {
	meta := PageMeta{}
	if meta.Clone().RemainingItemCount != nil {
		t.Fatal("nil RemainingItemCount clone is non-nil")
	}
}
