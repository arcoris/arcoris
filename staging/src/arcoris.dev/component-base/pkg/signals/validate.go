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
)

const (
	// errNilValidationMessage is the stable diagnostic text used when a
	// validation helper is called without a diagnostic message.
	//
	// Validation helpers are package-local and should always be called with the
	// stable diagnostic text owned by the feature file that rejects the input.
	errNilValidationMessage = "signals: nil validation message"
)

// requireValidationMessage panics when message is empty.
func requireValidationMessage(message string) {
	if message == "" {
		panic(errNilValidationMessage)
	}
}

// requireContext panics when ctx is nil.
func requireContext(ctx context.Context, message string) {
	requireValidationMessage(message)
	if ctx == nil {
		panic(message)
	}
}

// requireSignal panics when sig is nil.
func requireSignal(sig os.Signal, message string) {
	requireValidationMessage(message)
	if sig == nil {
		panic(message)
	}
}

// requireSignals panics when any signal in sigs is nil.
func requireSignals(sigs []os.Signal, message string) {
	requireValidationMessage(message)
	for _, sig := range sigs {
		if sig == nil {
			panic(message)
		}
	}
}

// requireNonEmptySignals panics when sigs is empty or contains nil signals.
func requireNonEmptySignals(sigs []os.Signal, message string) {
	requireValidationMessage(message)
	if len(sigs) == 0 {
		panic(message)
	}
	requireSignals(sigs, message)
}

// requirePositiveBuffer panics when size is not positive.
func requirePositiveBuffer(size int, message string) {
	requireValidationMessage(message)
	if size <= 0 {
		panic(message)
	}
}

// requireNonNegativeBuffer panics when size is negative.
func requireNonNegativeBuffer(size int, message string) {
	requireValidationMessage(message)
	if size < 0 {
		panic(message)
	}
}
