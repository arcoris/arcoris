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

import "errors"

// WaitError describes a failed lifecycle wait operation.
//
// WaitError carries the snapshot observed when the wait failed. For state-based
// waits, Target identifies the desired state. The wrapped Err identifies why the
// wait failed, such as ErrWaitTargetUnreachable, ErrInvalidWaitPredicate,
// ErrInvalidWaitTarget, context.Canceled, or context.DeadlineExceeded.
//
// WaitError is intended for wait.go. It should not be used for transition
// validation or guard rejection.
type WaitError struct {
	// Snapshot is the lifecycle snapshot observed when the wait failed.
	//
	// For ErrWaitTargetUnreachable, Snapshot is the state that made the target
	// impossible to reach. For context cancellation or deadline expiration,
	// Snapshot is the most recent lifecycle observation available to the wait
	// implementation.
	Snapshot Snapshot

	// Target is the state the caller was waiting for.
	//
	// Target is meaningful only when HasTarget is true. Generic predicate waits
	// should leave HasTarget false because StateNew is a valid zero-value state
	// and cannot be used as an implicit "no target" marker.
	Target State

	// HasTarget reports whether Target is meaningful.
	HasTarget bool

	// Err is the reason the wait failed.
	//
	// Err SHOULD be one of:
	//
	//   - ErrWaitTargetUnreachable;
	//   - ErrInvalidWaitPredicate;
	//   - ErrInvalidWaitTarget;
	//   - context.Canceled;
	//   - context.DeadlineExceeded.
	//
	// If Err is nil, the error is treated as ErrWaitTargetUnreachable.
	Err error
}

// Error returns a Go-style diagnostic message for the failed wait.
func (e *WaitError) Error() string {
	cause := waitErrorCause(e)

	if e == nil {
		return cause.Error()
	}

	message := waitErrorMessagePrefix(cause)

	if e.HasTarget {
		return message +
			": target " +
			e.Target.String() +
			" at " +
			e.Snapshot.String()
	}

	return message + ": at " + e.Snapshot.String()
}

// Unwrap returns the wait failure cause.
//
// This preserves context cancellation and deadline semantics:
//
//	errors.Is(err, context.Canceled)
//	errors.Is(err, context.DeadlineExceeded)
//
// It also preserves lifecycle wait semantics:
//
//	errors.Is(err, lifecycle.ErrWaitTargetUnreachable)
//	errors.Is(err, lifecycle.ErrInvalidWaitPredicate)
//	errors.Is(err, lifecycle.ErrInvalidWaitTarget)
func (e *WaitError) Unwrap() error {
	return waitErrorCause(e)
}

// newWaitError constructs a WaitError for a generic predicate wait.
//
// The helper is package-local because wait.go should be the normal producer of
// wait errors. Public callers can still inspect WaitError values returned by
// wait methods.
func newWaitError(snap Snapshot, err error) *WaitError {
	if err == nil {
		err = ErrWaitTargetUnreachable
	}

	return &WaitError{
		Snapshot: snap,
		Err:      err,
	}
}

// newWaitStateError constructs a WaitError for a state-targeted wait.
//
// HasTarget is set explicitly because StateNew is a valid zero-value state and
// cannot be used to represent an absent target.
func newWaitStateError(snap Snapshot, target State, err error) *WaitError {
	if err == nil {
		err = ErrWaitTargetUnreachable
	}

	return &WaitError{
		Snapshot:  snap,
		Target:    target,
		HasTarget: true,
		Err:       err,
	}
}

// waitErrorCause returns the effective cause of e.
func waitErrorCause(e *WaitError) error {
	if e == nil || e.Err == nil {
		return ErrWaitTargetUnreachable
	}

	return e.Err
}

// waitErrorMessagePrefix returns the stable message prefix for a wait error.
//
// Lifecycle wait sentinels already carry the "lifecycle:" package prefix. Other
// errors, such as context.Canceled or context.DeadlineExceeded, are wrapped in a
// lifecycle wait message so the final Error string still starts with
// "lifecycle:".
func waitErrorMessagePrefix(cause error) string {
	switch {
	case errors.Is(cause, ErrWaitTargetUnreachable):
		return ErrWaitTargetUnreachable.Error()
	case errors.Is(cause, ErrInvalidWaitPredicate):
		return ErrInvalidWaitPredicate.Error()
	case errors.Is(cause, ErrInvalidWaitTarget):
		return ErrInvalidWaitTarget.Error()
	default:
		return "lifecycle: wait failed: " + cause.Error()
	}
}
