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

package objectstorewatch

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectwatch"
)

func TestErrorSentinels(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "nil backend", err: ErrNilBackend},
		{name: "invalid option", err: ErrInvalidOption},
		{name: "stream overflow", err: ErrStreamOverflow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !errors.Is(tt.err, tt.err) {
				t.Fatalf("errors.Is(%v, %v) = false", tt.err, tt.err)
			}
		})
	}
}

func TestStreamOverflowErrorPreservesContinuityAndOverflow(t *testing.T) {
	err := streamOverflowError()

	if !errors.Is(err, objectwatch.ErrContinuityLost) {
		t.Fatalf("errors.Is(%v, %v) = false", err, objectwatch.ErrContinuityLost)
	}
	if !errors.Is(err, ErrStreamOverflow) {
		t.Fatalf("errors.Is(%v, %v) = false", err, ErrStreamOverflow)
	}
}
