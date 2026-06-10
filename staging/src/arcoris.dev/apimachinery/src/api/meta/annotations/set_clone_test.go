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

package annotations

import "testing"

func TestSetClone(t *testing.T) {
	if Set(nil).Clone() != nil {
		t.Fatal("nil clone is non-nil")
	}

	empty := Set{}
	emptyClone := empty.Clone()
	if emptyClone == nil {
		t.Fatal("empty clone is nil")
	}
	emptyClone["note"] = "created"
	if len(empty) != 0 {
		t.Fatal("empty clone aliases original map")
	}

	set := Set{"note": "human readable"}
	cloned := set.Clone()
	cloned["note"] = "changed"

	if set["note"] != "human readable" {
		t.Fatal("Clone() did not detach map storage")
	}
	set["tool"] = "cli"
	if _, ok := cloned["tool"]; ok {
		t.Fatal("original mutation changed clone")
	}
}
