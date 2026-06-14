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

package objectsurface

import "testing"

func TestKnownKindsRoundTripThroughObjectPaths(t *testing.T) {
	for _, kind := range knownKinds() {
		t.Run(kind.String(), func(t *testing.T) {
			path, ok := kind.ObjectPath()
			if !ok {
				t.Fatalf("ObjectPath() ok = false")
			}

			roundTripped, ok := KindFromObjectPath(path)
			if !ok {
				t.Fatalf("KindFromObjectPath(%q) ok = false", path)
			}
			if roundTripped != kind {
				t.Fatalf("KindFromObjectPath(%q) = %q; want %q", path, roundTripped, kind)
			}
			if roundTripped.String() == path.String() {
				t.Fatalf("surface ID %q unexpectedly equals full object path", roundTripped)
			}
		})
	}
}
