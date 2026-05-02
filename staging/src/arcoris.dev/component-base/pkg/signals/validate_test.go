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

package signals

import (
	"context"
	"os"
	"testing"
)

func TestValidationHelpersRejectMissingDiagnosticText(t *testing.T) {
	mustPanicWith(t, errNilValidationMessage, func() {
		requireContext(context.Background(), "")
	})
}

func TestValidationHelpersRejectInvalidInputs(t *testing.T) {
	mustPanicWith(t, "nil context", func() {
		requireContext(nil, "nil context")
	})
	mustPanicWith(t, "nil signal", func() {
		requireSignal(nil, "nil signal")
	})
	mustPanicWith(t, "nil signals", func() {
		requireSignals([]os.Signal{nil}, "nil signals")
	})
	mustPanicWith(t, "empty signals", func() {
		requireNonEmptySignals(nil, "empty signals")
	})
	mustPanicWith(t, "positive buffer", func() {
		requirePositiveBuffer(0, "positive buffer")
	})
	mustPanicWith(t, "non-negative buffer", func() {
		requireNonNegativeBuffer(-1, "non-negative buffer")
	})
}

func TestNoCopyMarkerMethodsAreCallable(t *testing.T) {
	var marker noCopy
	marker.Lock()
	marker.Unlock()
}
