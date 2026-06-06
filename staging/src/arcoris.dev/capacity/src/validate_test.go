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

package capacity

import (
	"errors"
	"testing"
)

func TestRequireValidVectorPanicsOnInvalidVector(t *testing.T) {
	requirePanicMatches(t, ErrInvalidVector, func() {
		requireValidVector("vector", Vector{entries: []Entry{{Resource: "bad-name", Amount: 1}}})
	})
}

func TestRequireValidDemandPanicsOnInvalidDemand(t *testing.T) {
	requirePanicMatches(t, ErrInvalidDemand, func() {
		requireValidDemand("demand", Demand{})
	})
}

func TestRequirePositiveAmountPanicsOnZero(t *testing.T) {
	requirePanicMatches(t, ErrZeroAmount, func() {
		requirePositiveAmount(0)
	})
}

func requirePanicMatches(t *testing.T, sentinel error, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("panic = nil, want %v", sentinel)
		}
		if !errors.Is(recovered.(error), sentinel) {
			t.Fatalf("panic = %v, want %v", recovered, sentinel)
		}
	}()

	fn()
}
