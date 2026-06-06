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

func TestErrorAtBuildsCapacityError(t *testing.T) {
	err := errorAt("amount", ErrZeroAmount, "amount must be positive")

	var capacityErr *Error
	if !errors.As(err, &capacityErr) {
		t.Fatalf("errorAt() = %T, want *Error", err)
	}
	if capacityErr.Path != "amount" || !errors.Is(capacityErr, ErrZeroAmount) {
		t.Fatalf("capacity error = %#v", capacityErr)
	}
}

func TestPanicAtPanicsWithCapacityError(t *testing.T) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("panicAt() did not panic")
		}
		if !errors.Is(recovered.(error), ErrZeroAmount) {
			t.Fatalf("panic error = %v, want ErrZeroAmount", recovered)
		}
	}()

	panicAt("amount", ErrZeroAmount, "amount must be positive")
}
