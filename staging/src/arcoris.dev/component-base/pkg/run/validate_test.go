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

func TestValidationHelpersRejectMissingDiagnosticText(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilValidationMessage, func() {
		requireContext(context.Background(), "")
	})
}

func TestValidationHelpersRejectInvalidInputs(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, "nil context", func() {
		requireContext(nil, "nil context")
	})
	mustPanicWith(t, errNilGroup, func() {
		requireGroup(nil)
	})
	mustPanicWith(t, errUninitializedGroup, func() {
		requireGroup(&Group{})
	})
	mustPanicWith(t, errNilTask, func() {
		requireTask(nil)
	})
	mustPanicWith(t, errEmptyTaskName, func() {
		requireTaskName("")
	})
	mustPanicWith(t, errUntrimmedTaskName, func() {
		requireTaskName(" worker")
	})
	mustPanicWith(t, errInvalidErrorMode, func() {
		requireErrorMode(ErrorMode(99))
	})
}

func TestNoCopyMarkerMethodsAreCallable(t *testing.T) {
	t.Parallel()

	var marker noCopy
	marker.Lock()
	marker.Unlock()
}
