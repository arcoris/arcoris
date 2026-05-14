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

package run

import (
	"context"
	"testing"
)

func TestTaskExecutesWithContext(t *testing.T) {
	t.Parallel()

	called := false
	task := Task(func(ctx context.Context) error {
		called = true
		if ctx == nil {
			t.Fatal("task received nil context")
		}
		return nil
	})

	if err := task(context.Background()); err != nil {
		t.Fatalf("task error = %v", err)
	}
	if !called {
		t.Fatal("task was not called")
	}
}

func TestTaskZeroValueIsNil(t *testing.T) {
	t.Parallel()

	var task Task
	if task != nil {
		t.Fatal("zero Task is non-nil")
	}
}
