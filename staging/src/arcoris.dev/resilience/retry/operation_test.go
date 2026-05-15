/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package retry

import (
	"context"
	"testing"
)

// TestOperationAcceptsContextAwareFunction documents that Operation is the
// function shape consumed by Do. The test is intentionally compile-oriented
// because Operation itself has no runtime behavior.
func TestOperationAcceptsContextAwareFunction(t *testing.T) {
	var op Operation = func(context.Context) error {
		return nil
	}

	if err := op(context.Background()); err != nil {
		t.Fatalf("Operation returned unexpected error: %v", err)
	}
}

// TestValueOperationAcceptsContextAwareFunction documents that ValueOperation is
// the value-returning function shape consumed by DoValue. The test is
// intentionally compile-oriented because ValueOperation itself has no runtime
// behavior.
func TestValueOperationAcceptsContextAwareFunction(t *testing.T) {
	var op ValueOperation[int] = func(context.Context) (int, error) {
		return 42, nil
	}

	val, err := op(context.Background())
	if err != nil {
		t.Fatalf("ValueOperation returned unexpected error: %v", err)
	}
	if val != 42 {
		t.Fatalf("ValueOperation returned %d, want 42", val)
	}
}
