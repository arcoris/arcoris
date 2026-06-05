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

// requireNonNil panics when l is a nil receiver.
func (l *Ledger) requireNonNil() {
	if l == nil {
		panicAt("ledger", ErrNilLedger, ErrorReasonNilLedger, "ledger receiver is nil")
	}
}

// requireInitializedLocked panics when l is a zero-value Ledger.
//
// The caller must hold l.mu so revision is read from a stable owner state.
func (l *Ledger) requireInitializedLocked() {
	if l.revision.IsZero() {
		panicAt(
			"ledger",
			ErrUninitializedLedger,
			ErrorReasonUninitializedLedger,
			"ledger must be created with NewLedger",
		)
	}
}
