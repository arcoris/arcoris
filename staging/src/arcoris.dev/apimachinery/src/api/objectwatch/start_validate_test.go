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

func TestStartValidate(t *testing.T) {
	tests := []struct {
		name  string
		start Start
		valid bool
	}{
		{name: "after revision", start: Start{Mode: StartAfterRevision, Revision: 1}, valid: true},
		{name: "after zero", start: Start{Mode: StartAfterRevision}, valid: false},
		{name: "current", start: Start{Mode: StartAtCurrent}, valid: true},
		{name: "current non-zero", start: Start{Mode: StartAtCurrent, Revision: 1}, valid: false},
		{name: "unknown mode", start: Start{Mode: StartMode(99)}, valid: false},
		{name: "zero", start: Start{}, valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.start.Validate()
			if tt.valid {
				requireNoError(t, err)
				return
			}
			requireErrorIs(t, err, ErrInvalidStart)
		})
	}
}

func TestStartValidateRejectsZeroRevisionType(t *testing.T) {
	err := (Start{Mode: StartAfterRevision, Revision: objectstore.Revision(0)}).Validate()

	requireErrorIs(t, err, ErrInvalidStart)
}
