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

package objectwatch

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestAfterRevision(t *testing.T) {
	tests := []struct {
		name     string
		revision int
	}{
		{name: "zero", revision: 0},
		{name: "non-zero", revision: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, err := AfterRevision(objectstore.Revision(tt.revision))
			requireNoError(t, err)

			if start.Mode != StartAfterRevision || start.Revision != objectstore.Revision(tt.revision) {
				t.Fatalf("start = %#v; want after revision %d", start, tt.revision)
			}
			if !start.IsValid() || start.IsZero() {
				t.Fatalf("valid start flags = valid %v zero %v; want true/false", start.IsValid(), start.IsZero())
			}
		})
	}
}

func TestAtCurrent(t *testing.T) {
	start := AtCurrent()

	if start.Mode != StartAtCurrent || !start.Revision.IsZero() {
		t.Fatalf("start = %#v; want at current", start)
	}
	requireNoError(t, start.Validate())
}
