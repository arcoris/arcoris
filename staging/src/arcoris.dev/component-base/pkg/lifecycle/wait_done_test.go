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

package lifecycle

import (
	"errors"
	"testing"
)

func TestDoneReturnsStableChannel(t *testing.T) {
	t.Parallel()

	controller := NewController()
	first := controller.Done()
	second := controller.Done()
	if first != second {
		t.Fatal("Done returned different channels")
	}
}

func TestDoneOpenBeforeTerminalState(t *testing.T) {
	t.Parallel()

	// Done is signal-only: it remains open through active non-terminal states and
	// carries no state payload of its own.
	controller := NewController()
	mustNotSignalClosed(t, controller.Done())
	_, _ = controller.BeginStart()
	mustNotSignalClosed(t, controller.Done())
}

func TestDoneClosesAfterStopped(t *testing.T) {
	t.Parallel()

	controller := NewController()
	done := controller.Done()
	_, _ = controller.BeginStop()
	mustSignalClosed(t, done)
}

func TestDoneClosesAfterFailed(t *testing.T) {
	t.Parallel()

	controller := NewController()
	_, _ = controller.BeginStart()
	done := controller.Done()
	_, _ = controller.MarkFailed(errors.New("failed"))
	mustSignalClosed(t, done)
}

func TestDoneWorksOnZeroValueController(t *testing.T) {
	t.Parallel()

	var controller Controller
	done := controller.Done()
	if done == nil {
		t.Fatal("Done returned nil for zero-value Controller")
	}
	mustNotSignalClosed(t, done)
	_, _ = controller.BeginStop()
	mustSignalClosed(t, done)
}
