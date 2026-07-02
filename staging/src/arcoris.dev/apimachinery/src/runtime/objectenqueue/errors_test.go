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

package objectenqueue

import (
	"errors"
	"testing"
)

func TestErrorsAreDistinctSentinels(t *testing.T) {
	sentinels := []error{
		ErrNilQueue,
		ErrNilMapper,
		ErrNilEmit,
		ErrInvalidHandler,
	}

	for i, left := range sentinels {
		for j, right := range sentinels {
			if i == j {
				requireErrorIs(t, left, right)
				continue
			}
			if errors.Is(left, right) {
				t.Fatalf("errors.Is(%v, %v) = true; want false", left, right)
			}
		}
	}
}
