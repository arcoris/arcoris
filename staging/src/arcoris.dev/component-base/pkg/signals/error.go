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
	"errors"
	"os"
)

const (
	// errNilSignalErrorSignal is the stable diagnostic text used when
	// NewSignalError receives a nil signal.
	//
	// SignalError represents a concrete process signal. Constructing one without
	// a signal would hide an invalid signal boundary and make context causes
	// ambiguous.
	errNilSignalErrorSignal = "signals: nil signal error signal"

	// errNilCauseContext is the stable diagnostic text used when Cause receives a
	// nil context.
	//
	// Cause is an inspection helper around context.Cause. A nil context has no
	// cancellation state to inspect and indicates invalid caller code.
	errNilCauseContext = "signals: nil cause context"
)

var (
	// ErrSignal is matched by SignalError values created by this package.
	//
	// Use errors.Is(err, ErrSignal) to classify a cancellation cause as signal
	// owned without depending on the concrete signal value.
	ErrSignal = errors.New("signal received")

	// ErrStopped reports that an owner-controlled signal Subscription was stopped
	// before Wait received a signal.
	//
	// ErrStopped is not a signal cause. It represents owner cleanup, not process
	// signal delivery.
	ErrStopped = errors.New("signal subscription stopped")
)

// SignalError is the typed context cancellation cause used for received signals.
//
// SignalError matches ErrSignal with errors.Is and exposes the Event that caused
// cancellation. It is suitable for use with context.WithCancelCause and can be
// recovered later through errors.As or Cause.
type SignalError struct {
	// Event is the signal event that caused cancellation.
	Event Event
}

// NewSignalError returns a SignalError for sig.
//
// NewSignalError panics when sig is nil.
func NewSignalError(sig os.Signal) SignalError {
	requireSignal(sig, errNilSignalErrorSignal)

	return SignalError{Event: Event{Signal: sig}}
}

// Error returns a stable human-readable signal error message.
func (e SignalError) Error() string {
	if e.Event.Signal == nil {
		return "signal received: <nil>"
	}

	return "signal received: " + e.Event.Signal.String()
}

// Is reports whether target is ErrSignal.
func (e SignalError) Is(target error) bool {
	return target == ErrSignal
}

// Cause extracts a signal Event from ctx's cancellation cause.
//
// Cause returns false when the context has not been cancelled by this package or
// when the cancellation cause is not a SignalError. Cause panics when ctx is nil.
func Cause(ctx context.Context) (Event, bool) {
	requireContext(ctx, errNilCauseContext)

	var signalErr SignalError
	if errors.As(context.Cause(ctx), &signalErr) {
		return signalErr.Event, signalErr.Event.Signal != nil
	}

	return Event{}, false
}
